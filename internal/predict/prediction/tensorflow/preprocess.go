package tensorflow

import (
	"github.com/pkg/errors"
	tf "github.com/wamuir/graft/tensorflow"
	"github.com/wamuir/graft/tensorflow/op"
)

const (
	errorTextCouldNotCreateTensorFromImage     = "could not create tensor from input image"
	errorTextCouldNotDecodeJPEG                = "could not decode jpeg channels of image"
	errorTextCouldNotRunPreprocessImageSession = "could not run tensorflow session to preprocess input image"
	int32zero                                  = int32(0)
	vgg16ImagenetMeanRed                       = float32(123.68)
	vgg16ImagenetMeanGreen                     = float32(116.779)
	vgg16ImagenetMeanBlue                      = float32(103.939)
)

// VGG16 mean RGB values for the imagenet dataset
var imagenetMeans = []float32{vgg16ImagenetMeanRed, vgg16ImagenetMeanGreen, vgg16ImagenetMeanBlue}

// Preprocessing in specific to VGG16
func (s *Service) makeTensorFromImage(imageBytes []byte, colorChannels int64) (*tf.Tensor, error) {

	// DecodeJpeg uses a scalar String-valued tensor as inputOperation.
	tensor, err := tf.NewTensor(string(imageBytes))
	if err != nil {
		return nil, errors.Wrap(err, errorTextCouldNotCreateTensorFromImage)
	}
	// Creates a tensorflow graph to decode the jpeg image
	graph, input, output, err := decodeJPEGGraph(colorChannels)
	if err != nil {
		return nil, errors.Wrap(err, errorTextCouldNotDecodeJPEG)
	}
	// Execute that graph to decode this one image
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	normalized, err := session.Run(
		map[tf.Output]*tf.Tensor{input: tensor},
		[]tf.Output{output},
		nil)
	if err != nil {
		return nil, errors.Wrap(err, errorTextCouldNotRunPreprocessImageSession)
	}
	return normalized[0], nil
}

func decodeJPEGGraph(colorChannels int64) (graph *tf.Graph, input, output tf.Output, err error) {
	s := op.NewScope()

	mean := op.Const(s, imagenetMeans)
	input = op.Placeholder(s, tf.String)
	output = op.DecodeJpeg(s, input, op.DecodeJpegChannels(colorChannels))
	output = op.Cast(s, output, tf.Float)
	output = op.Sub(s, output, mean)
	output = op.ExpandDims(s, output, op.Const(s.SubScope("batch"), int32zero))

	graph, err = s.Finalize()
	return graph, input, output, err
}
