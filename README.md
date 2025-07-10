# LiveTran

A simple, self-hostable live streaming server written in Go. LiveTran ingests video streams via the SRT protocol, transcodes them into HLS in real-time, and makes them available for playback. It also supports automatically uploading HLS segments to Cloudflare R2 for scalable delivery.

## Features

- **SRT Ingest:** Secure and reliable video ingest using the SRT protocol.
- **Real-time Transcoding:** Uses FFmpeg to transcode incoming streams into HLS (HTTP Live Streaming) format.
- **Stream Management API:** A simple JSON-based HTTP API to start, stop, and check the status of streams.
- **Live Playback:** Serves HLS playlists and segments directly for live viewing.
- **Cloudflare R2 Integration:** Automatically watches and uploads HLS files to a Cloudflare R2 bucket.
- **Webhook Notifications:** Sends status updates (e.g., "STREAMING", "STOPPED") to configured webhook URLs.
- **Concurrent Streams:** Manages multiple independent streams simultaneously.

## Getting Started

### Prerequisites

- [Go](https://go.dev/doc/install) (version 1.21 or later)
- [FFmpeg](https://ffmpeg.org/download.html) must be installed and available in your system's `PATH`.

### Installation

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/your-username/LiveTran.git
    cd LiveTran
    ```

2.  **Configure your environment:**
    Create a `.env` file in the root of the project. To make this easier, you can copy the provided example file:
    ```sh
    cp .env.example .env
    ```
    Now, edit the `.env` file with your Cloudflare R2 credentials if you plan to use the upload feature.

3.  **Build and run the server:**
    Assuming your main package is in a `cmd/server` directory:
    ```sh
    go build -o livtran_server ./cmd/server
    ./livtran_server
    ```
    The server will start, by default on port `8080`.

## How It Works

1.  **Start a Stream:** You send a `POST` request to the `/stream/start` endpoint with a unique `stream_id`.
2.  **Get SRT URL:** The server allocates a free port, creates a unique SRT listener, and logs the SRT URL. This status is also sent to any configured webhooks. The URL looks like `srt://<server_ip>:<port>?streamid=<your_stream_id>`.
3.  **Stream from OBS:** You configure your streaming software (like OBS) to push the stream to this SRT URL. The `streamid` in the URL must match the one you provided in the API call.
4.  **Transcoding & Playback:** Once the stream is received, LiveTran pipes it to FFmpeg, which generates HLS files (`.m3u8` and `.ts`) in the `output/` directory.
5.  **Cloud Upload (Optional):** If R2 credentials are provided in the `.env` file, a file watcher automatically uploads these HLS files to your specified R2 bucket.
6.  **Viewing:** You can play the live stream using an HLS-compatible player (like VLC or a web player) with the URL: `http://<server_ip>:8080/output/<your_stream_id>.m3u8`.

## API Reference

The API server runs on port `8080`. All request and response bodies are in JSON format.

### Start a Stream

Starts a new transcoding task.

- **Endpoint:** `POST /stream/start`
- **Request Body:**
  ```json
  {
    "stream_id": "my-unique-stream-123",
    "webhook_urls": ["https://my-service.com/webhook"]
  }
  ```
- **Success Response:**
  ```json
  {
    "success": true,
    "data": "Stream launching!"
  }
  ```

### Stop a Stream

Stops an active stream.

- **Endpoint:** `POST /stream/stop`
- **Request Body:**
  ```json
  {
    "stream_id": "my-unique-stream-123"
  }
  ```
- **Success Response:**
  ```json
  {
    "success": true,
    "data": "Stream stopped!"
  }
  ```

### Get Stream Status

Checks the current status of a stream.

- **Endpoint:** `POST /stream/status`
- **Request Body:**
  ```json
  {
    "stream_id": "my-unique-stream-123"
  }
  ```
- **Success Response:**
  ```json
  {
    "success": true,
    "data": "Status: STREAMING"
  }
  ```
  Possible statuses are: `INITIALISED`, `READY`, `STREAMING`, `STOPPED`.

### HLS Playback

Serves the HLS files for playback directly from the server.

- **Playlist:** `GET /output/{stream_id}.m3u8`
- **Segments:** `GET /output/{stream_id}-*.ts`

## Future Work

Based on the `TODO`s in the codebase, here are some planned improvements:

- **Authentication:** Implement JWT and Stream Key validation for securing the API and ingest endpoints.
- **Resource Cleanup:** Automatically delete local HLS files after they have been successfully uploaded to R2.
- **Improved File Handling:** Replace `time.Sleep` with a more robust mechanism to ensure files are complete before uploading.

## License

This project is licensed under the MIT License. See the LICENSE file for details.