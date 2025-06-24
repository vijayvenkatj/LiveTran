package ingest

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/datarhei/gosrt"
)


func SrtConnectionTask(ctx context.Context,task *Task) {

	port,err := getFreePort()
	if err != nil {
		task.UpdateStatus(StreamStopped,fmt.Sprintf("PORT error: %s",err))
		return
	}

	addr := fmt.Sprintf(":%d",port)
	listener, err := srt.Listen("srt", addr, srt.DefaultConfig())
	if err != nil {
		task.UpdateStatus(StreamStopped,fmt.Sprintf("SRT Listener error: %s",err))
		return
	}

	defer listener.Close()

	ip := GetLocalIP()
	url := fmt.Sprintf("srt://%s:%d?streamid=%s",ip,port,task.Id)

	task.UpdateStatus(StreamReady,fmt.Sprintf("The stream is ready! URL -> %s",url))

	defer close(task.UpdatesChan)



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
		
		// User stops the stream.
		case <- ctx.Done():
			listener.Close()
			task.UpdateStatus(StreamStopped, fmt.Sprintf("Stream is stopped : %v", context.Cause(ctx)))
			return 

		// SRT connection timeout because of inactivity.
		case <- time.After(120*time.Second):
			listener.Close()
			task.UpdateStatus(StreamStopped,"TIMEOUT due to inactivity!")
			return
		
		// Error from the Accept2() routine
		case err := <- errChan:
			fmt.Println("ERROR: ",err)
			continue
		
		// Request gets detected by the Accept2() routine.
		case req := <- reqChan:

			if req.StreamId() != task.Id {
				req.Reject(srt.REJ_BADSECRET)
				fmt.Println("SRT ERROR: StreamId does not match.")
				continue
			}

			conn,err := req.Accept()
			if err != nil {
				fmt.Println("SRT Request Error : ",err)
				return
			}

			task.UpdateStatus(StreamActive,"Your stream is now live!")
			handleStream(ctx,conn,task)
		}
	}

}


func handleStream(ctx context.Context,conn srt.Conn,task *Task) {
	defer conn.Close()
	buf := make([]byte, 1316) // 1316 bytes is a typical MPEG-TS packet size

	go func(){
		<- ctx.Done()
		conn.Close()
	}()

	for {
		n, err := conn.Read(buf)
		if err != nil {
			task.UpdateStatus(StreamReady,fmt.Sprintf("Connection closed : %s",err))
			return
		}

		// TODO: process the incoming stream data (e.g., forward to FFmpeg)
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
