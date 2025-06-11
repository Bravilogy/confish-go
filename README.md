# Confish Go Client

A lightweight Go client for interacting with the [confish](https://confi.sh) configuration platform. This client allows you to fetch configuration values, send logs, and handle webhook updates in your Go applications.

## ✨ Features

- Fetch dynamic configurations by ID
- Send structured logs to the Confish logging API
- Process webhook payloads from confish
- Built-in log level helpers (`Debug`, `Info`, `Warn`, etc.)

---

## 📦 Installation

```bash
go get github.com/bravilogy/confish-go@latest
```

---

## 🛠 Usage

### 1. Import the package

```go
import "github.com/bravilogy/confish-go/confish"
```

### 2. Initialize the client

```go
cfg := &confish.ConfishConfig{
    URL:       "https://api.confi.sh",
    AppID:     "your-app-id",
    AppSecret: "your-app-secret",
}

client, err := confish.NewClient(cfg)
if err != nil {
    log.Fatalf("failed to create client: %v", err)
}
```

### 3. Fetch a configuration

```go
var configStruct struct {
    FeatureEnabled bool   `json:"feature_enabled"`
    TimeoutSeconds int    `json:"timeout_seconds"`
    Environment    string `json:"environment"`
}

err = client.GetConfig("your-config-id", &configStruct)
if err != nil {
    log.Fatalf("failed to fetch config: %v", err)
}

fmt.Printf("Fetched config: %+v\n", configStruct)
```

### 4. Send a log message

```go
err = client.Info("Starting background worker")
if err != nil {
    log.Printf("failed to log info: %v", err)
}
```

You can also use other levels:

```go
client.Debug("Detailed debug log")
client.Warn("Something suspicious")
client.Error("Something failed")
client.Critical("System is down")
```

### 5. Handle a webhook payload

Assuming you have an HTTP handler set up for your webhook call:

```go
var payload confish.WebhookPayload
err := json.NewDecoder(req.Body).Decode(&payload)
if err != nil {
    http.Error(w, "invalid payload", http.StatusBadRequest)
    return
}

var updatedValues struct {
    FeatureToggle bool `json:"feature_toggle"`
}

err = client.ProcessWebhookPayload(payload, &updatedValues)
if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
}

fmt.Printf("Updated values: %+v\n", updatedValues)
```

---

## 🔐 Authentication

Every request requires:

- `App-ID` – Your unique application identifier
- `App-Secret` – Your secret token to authenticate

These should be passed to the client via the `ConfishConfig` struct.

---

## 📑 License

MIT

---

## 📬 Support

For issues or feature requests, please [open an issue](https://github.com/bravilogy/confish-go/issues).

To learn more about the confish platform, visit [https://confi.sh](https://confi.sh).
