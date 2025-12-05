package main

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"sync"
	"time"
)

const (
	KUBECTL = "kubctl"
	GIT     = "git"
	PING    = "ping"
	SSH     = "ssh"
)

type sshRecord struct {
	Host      net.IP
	Reachable bool
	LoginSSH  bool
	Uname     string
}

func runUname(ctx context.Context, host net.IP, user string) (string, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}
	login := fmt.Sprintf("%s@%s", user, host)
	cmd := exec.CommandContext(ctx, SSH, "-o StrictHostKeyChecking=no", "-o BatchMode=yes", login, "uname -a")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func hostAlive(ctx context.Context, host net.IP) bool {
	cmd := exec.CommandContext(ctx, PING, "-c", "1", "-t", "2", host.String())
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func scanPrefixes(ipCh chan net.IP) chan sshRecord {
	ch := make(chan sshRecord, 1)
	go func()  {
		defer close(ch)
		limit := make(chan struct{}, 100)
		wg := sync.WaitGroup{}
		for ip := range ipCh {
			limit <- struct{}{}
			wg.Add(1)
			go func (ip net.IP)  {
				defer func ()  {
					<-limit
				}()
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
				defer cancel()
				rec := sshRecord{
					Host: ip,
				}
				if hostAlive(ctx, ip) {
					rec.Reachable = true
				}
				ch <- rec
			}(ip)
		}
		wg.Wait()
	}()
	return ch
}
