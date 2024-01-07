package gateway

import "testing"

func TestGatewayService_Init(t *testing.T) {
	service := NewDefaultGatewayService()
	service.start()
}
