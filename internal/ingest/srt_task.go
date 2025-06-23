package ingest

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/datarhei/gosrt"
)


func SrtConnectionTask(ctx context.Context,streamId string) error {

	port,err := getFreePort()
	if err != nil {
		fmt.Println("PORT error:",err)
		return err
	}

	addr := fmt.Sprintf(":%d",port)
	listener, err := srt.Listen("srt", addr, srt.DefaultConfig())
	if err != nil {
		fmt.Println("SRT Listener error:", err)
		return err
	}

	defer listener.Close()

	ip := GetLocalIP()
	url := fmt.Sprintf("srt://%s:%d?streamid=%s",ip,port,streamId)
	fmt.Println("OBS URL : ", url)

	/*
	PROBLEM:
		The Accept2() function is working always to catch the incoming stream.
		Suppose if i use this in this routine, then I can't control this properly because Default would block the events. 
		ie ) Flow wont reach <- ctx.Done()
	SOLN: 
		Make a seperate goroutine to check for Accept2() 
			IF error then send the error to the error buffered channel
			IF Connection is found add it to reqChan
	*/

	errChan := make(chan error,1)
	reqChan := make(chan srt.ConnRequest,1)

	go func() {
		for {
			select {
			case <- ctx.Done():
				return 
			default:
				req,err := listener.Accept2()
				if err != nil {
					errChan <- err
					return 
				}
				reqChan <- req
			}
		}
	}()

	for {
		select {
		case <- ctx.Done():
			fmt.Println(context.Cause(ctx))
			listener.Close()
			return errors.New("srt connection closed by user")
		case err := <- errChan:
			fmt.Println("ERROR: ",err)
			continue
		case req := <- reqChan:
			if req.StreamId() != streamId {
				req.Reject(srt.REJ_BADSECRET)
				fmt.Println("SRT ERROR: StreamId does not match.")
				continue
			}

			conn,err := req.Accept()
			if err != nil {
				fmt.Println("SRT Request Error : ",err)
				return err
			}

			handleStream(conn)
		}
	}

}

/*	
	TODO: 
		After The pod gets an OBS url, We send it to the user using webhooks.
		We add a Status to each job. -> Initialising, Ready, Streaming, Stopped
*/

func handleStream(conn srt.Conn) {
	defer conn.Close()
	buf := make([]byte, 1316) // 1316 bytes is a typical MPEG-TS packet size

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Connection closed or errored:", err)
			return
		}

		// process the incoming stream data (e.g., forward to FFmpeg)
		fmt.Printf("Received %d bytes\n", n)
	}
}


func GetLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return "127.0.0.1"
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return "127.0.0.1"
}

func getFreePort() (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
