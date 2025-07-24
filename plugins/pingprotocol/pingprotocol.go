package pingprotocol

import (
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proto/version"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/gate/proto"
	"gopkg.in/yaml.v3"
)

type Protocol struct {
	Number *int     `yaml:"number"`
	Names  []string `yaml:"names"`
}

type Config struct {
	Protocols []Protocol `yaml:"protocols"`
}

var Plugin = proxy.Plugin{
	Name: "PingProtocol",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		log := logr.FromContextOrDiscard(ctx)

		data, err := os.ReadFile("protocol.yml")
		if err != nil {
			log.Error(err, "Unable to read config file")
			return err
		}

		var cfg Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			log.Error(err, "Unable to read config file")
			return err
		}
		supportedVersions := func() (vs []*proto.Version) {
			for _, p := range cfg.Protocols {
				if p.Number == nil {
					log.Info("Protocol number is missing, skipping", "protocol", p)
					continue
				}
				var (
					v  *proto.Version
					ok bool
				)
				if len(p.Names) != 0 {
					v = &proto.Version{Protocol: proto.Protocol(*p.Number), Names: p.Names}
				} else if v, ok = version.ProtocolToVersion[proto.Protocol(*p.Number)]; !ok {
					v = &proto.Version{Protocol: proto.Protocol(*p.Number), Names: []string{fmt.Sprintf("v%d", p.Number)}}
				}
				if v != nil {
					vs = append(vs, v)
				}
			}
			return
		}()

		event.Subscribe(p.Event(), 0, onPing(supportedVersions))

		log.Info("PingProtocol plugin init successfully!")

		return nil
	},
}

func onPing(supportedVersions []*proto.Version) func(*proxy.PingEvent) {
	supportedVersionsString := func() (s string) {
		var prevVersion int
		isContinuous := false
		for i, v := range supportedVersions {
			if i == 0 {
				prevVersion = int(v.Protocol)
				s += v.FirstName()
			} else {
				if prevVersion+1 != int(v.Protocol) {
					if isContinuous {
						s += "-" + supportedVersions[i-1].LastName()
						isContinuous = false
					}
					s += ", " + v.FirstName()
				} else {
					isContinuous = true
				}
				prevVersion = int(v.Protocol)
			}
		}
		lastVersion := supportedVersions[len(supportedVersions)-1]
		if isContinuous || len(lastVersion.Names) != 1 {
			s += "-" + lastVersion.LastName()
		}
		return
	}()
	return func(e *proxy.PingEvent) {
		clientVersion := version.Protocol(e.Connection().Protocol())

		p := e.Ping()
		p.Version.Name = supportedVersionsString

		// Check if client version is in supported protocols
		for _, v := range supportedVersions {
			if clientVersion == version.Protocol(v.Protocol) {
				p.Version.Protocol = v.Protocol
				return
			}
		}
		// Client version not supported, set protocol to first supported
		p.Version.Protocol = supportedVersions[0].Protocol
	}
}
