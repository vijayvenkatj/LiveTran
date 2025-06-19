# README.md

# LiveTran – Low-Latency Live Streaming Platform

LiveTran is an open-source, cloud-native live streaming platform inspired by Mux. It enables ultra-low latency (sub-5s) video streaming, adaptive bitrate transcoding, real-time analytics, and scalable infrastructure—all powered by Go, Kubernetes, and modern streaming protocols.

---

## 🚀 Features

- **Ultra-Low Latency**: Achieve 4–5s end-to-end latency with LL-HLS and WebRTC fallback
- **Adaptive Bitrate Streaming**: Automatic transcoding to multiple resolutions (1080p, 720p, 480p)
- **Secure Ingestion**: SRT/WebRTC ingest endpoints with JWT-secured stream keys
- **Real-Time Analytics**: Live viewer counts, engagement stats, and stream health metrics
- **Scalable by Design**: Kubernetes-native, auto-scaling transcoders, and hybrid CDN delivery
- **DevOps Ready**: Dockerized, with CI/CD pipeline, Prometheus/Grafana monitoring, and Terraform infrastructure

---

## 🏗️ Architecture Overview

```
Broadcaster → RTMP/WebRTC → [Go Media Server]
                   ↓
              [Transcoding]
                   ↓
      [LL-HLS/DASH Packaging] → [CDN Edge] → Viewer
```

---

## 📦 Tech Stack

- **Backend**: Go (Gin, gRPC), LiveGo, FFmpeg, Redis, PostgreSQL
- **Infrastructure**: Kubernetes, Docker, Terraform, Cloudflare CDN, Prometheus, Grafana

---

## 🌟 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/livepipe.git
cd livepipe
```

### 2. Setup Environment Variables

Copy `.env.example` to `.env` and fill in your secrets (DB, JWT, CDN keys, etc.).

### 3. Run Locally (Docker Compose)

```bash
docker-compose up --build
```

### 4. Access Services

- **API**: http://localhost:8080
- **Prometheus/Grafana**: http://localhost:9090 / :3001

### 5. Streaming

- Use OBS or any RTMP client to stream to:  
  `rtmp://localhost:1935/live/{your_stream_key}`

---

## 📖 Documentation

- [API Reference](docs/api.md)
- [Architecture](docs/architecture.md)
- [Deployment Guide](docs/deployment.md)
- [Streaming Setup](docs/streaming.md)

---

## 🛡️ Security

- JWT-secured endpoints
- Signed playback URLs
- Rate-limited API
- Encrypted HLS segments (AES-128)

---

## 📊 Analytics & Monitoring

- Real-time viewer and stream health dashboards (Prometheus/Grafana)
- Error and performance logs

---

## 🤝 Contributing

PRs and issues welcome! See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## 📜 License

MIT

---

## 🙌 Acknowledgements

- [LiveGo](https://github.com/gwuhaolin/livego)
- [FFmpeg](https://ffmpeg.org/)
- [Mux](https://mux.com/) for inspiration

---
