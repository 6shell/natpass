package pool

import (
	"natpass/code/client/global"
	"natpass/code/network"
	"time"

	"google.golang.org/protobuf/proto"
)

// SendConnectReq send connect request message
func (conn *Conn) SendConnectReq(id string, cfg global.Rule) {
	var msg network.Msg
	msg.To = cfg.Target
	msg.XType = network.Msg_connect_req
	msg.LinkId = id
	switch cfg.Type {
	case "shell":
		msg.Payload = &network.Msg_Creq{
			Creq: &network.ConnectRequest{
				Name:  cfg.Name,
				XType: network.ConnectRequest_shell,
				Payload: &network.ConnectRequest_Cshell{
					Cshell: &network.ConnectShell{
						Exec: cfg.Exec,
						Env:  cfg.Env,
					},
				},
			},
		}
	case "vnc":
		fps := cfg.Fps
		if fps > 50 {
			fps = 50
		} else if fps == 0 {
			fps = 10
		}
		msg.Payload = &network.Msg_Creq{
			Creq: &network.ConnectRequest{
				Name:  cfg.Name,
				XType: network.ConnectRequest_vnc,
				Payload: &network.ConnectRequest_Cvnc{
					Cvnc: &network.ConnectVnc{
						Fps: fps,
					},
				},
			},
		}
	}
	select {
	case conn.write <- &msg:
	case <-time.After(conn.parent.cfg.WriteTimeout):
	}
}

// SendConnectVnc send connect vnc request message
func (conn *Conn) SendConnectVnc(id string, cfg global.Rule, quality uint64, showCursor bool) {
	var msg network.Msg
	msg.To = cfg.Target
	msg.XType = network.Msg_connect_req
	msg.LinkId = id
	fps := cfg.Fps
	if fps > 50 {
		fps = 50
	} else if fps == 0 {
		fps = 10
	}
	msg.Payload = &network.Msg_Creq{
		Creq: &network.ConnectRequest{
			Name:  cfg.Name,
			XType: network.ConnectRequest_vnc,
			Payload: &network.ConnectRequest_Cvnc{
				Cvnc: &network.ConnectVnc{
					Fps:     fps,
					Quality: uint32(quality),
					Cursor:  showCursor,
				},
			},
		},
	}
	select {
	case conn.write <- &msg:
	case <-time.After(conn.parent.cfg.WriteTimeout):
	}
}

// SendConnectError send connect error response message
func (conn *Conn) SendConnectError(to string, toIdx uint32, id, info string) {
	var msg network.Msg
	msg.To = to
	msg.ToIdx = toIdx
	msg.XType = network.Msg_connect_rep
	msg.LinkId = id
	msg.Payload = &network.Msg_Crep{
		Crep: &network.ConnectResponse{
			Ok:  false,
			Msg: info,
		},
	}
	select {
	case conn.write <- &msg:
	case <-time.After(conn.parent.cfg.WriteTimeout):
	}
}

// SendConnectOK send connect success response message
func (conn *Conn) SendConnectOK(to string, toIdx uint32, id string) {
	var msg network.Msg
	msg.To = to
	msg.ToIdx = toIdx
	msg.XType = network.Msg_connect_rep
	msg.LinkId = id
	msg.Payload = &network.Msg_Crep{
		Crep: &network.ConnectResponse{
			Ok: true,
		},
	}
	select {
	case conn.write <- &msg:
	case <-time.After(conn.parent.cfg.WriteTimeout):
	}
}

// SendDisconnect send disconnect message
func (conn *Conn) SendDisconnect(to string, toIdx uint32, id string) uint64 {
	var msg network.Msg
	msg.To = to
	msg.ToIdx = toIdx
	msg.XType = network.Msg_disconnect
	msg.LinkId = id
	select {
	case conn.write <- &msg:
		data, _ := proto.Marshal(&msg)
		return uint64(len(data))
	case <-time.After(conn.parent.cfg.WriteTimeout):
		return 0
	}
}
