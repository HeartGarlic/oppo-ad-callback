package oppo_ad_callback

import (
	"testing"
)

func TestOppoAdCallback_SendData(t *testing.T) {
	config := OppoAdCallbackConfig{
		OwnerId: 0,
		ApiId:   "",
		ApiKey:  "",
	}
	oppo := NewOppoAdCallback(&config)
	params := SendDataParams{
		PageId:        0,
		Ip:            "",
		Tid:           "",
		LbId:          "",
		TransformType: TransFormTypeNewSubmit,
	}
	data, err := oppo.SendData(params)
	if err != nil {
		t.Errorf("SendData error: %v", err)
	}
	t.Logf("SendData response: %v", data)
}
