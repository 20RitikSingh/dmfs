# 🔐 Fast and Secure Chunk-Based File Hashing & Encryption

This project provides a high-performance, chunk-wise file processing pipeline for **hashing** and **encrypting** large files efficiently. Designed for server-side use, it leverages modern cryptographic primitives optimized for speed, security, and scalability.

---

## 📦 Overview

This system reads files in chunks, hashes them to generate a unique key, and encrypts the data efficiently. It is optimized for:

- ✅ Large file support (processed in memory-efficient chunks)
- ⚡ High-speed hashing
- 🔐 Secure, authenticated encryption
- 💻 Server-optimized performance

---

## ⚙️ Algorithms Used

### 🔍 Hashing: **BLAKE3**

BLAKE3 is a modern cryptographic hash function designed for speed and scalability.

| Feature            | Details |
|--------------------|---------|
| **Speed**          | 5–10× faster than SHA-256 |
| **Parallelizable** | Utilizes multi-core CPUs via SIMD |
| **Memory-Efficient** | Ideal for streaming large files in chunks |
| **Secure**         | Cryptographically safe (based on BLAKE2s) |
| **Output**         | 256-bit hash |

### ✅ Why BLAKE3?

- 🚀 **Fastest cryptographic hash** on modern CPUs
- 🌲 **Tree-hashing** structure — great for multi-threaded or streaming applications
- 💡 **Safe default**: Suitable for file integrity, fingerprints, deduplication

#### 🔄 Compared to Other Algorithms

| Algorithm     | Speed | Parallel | Secure | Notes |
|---------------|-------|----------|--------|-------|
| SHA-256       | ❌     | ❌        | ✅      | Slower, widely used |
| BLAKE2b       | ✅     | ⚠️ Partial | ✅      | Good but not as scalable |
| MD5 / SHA-1   | ✅     | ❌        | ❌      | Broken, avoid |
| **BLAKE3**    | ✅✅    | ✅        | ✅      | Best speed + security |

---

## 🔐 Encryption: **AES-GCM (Galois/Counter Mode)**

AES-GCM is a modern authenticated encryption mode built on AES-256.

| Feature              | Details |
|----------------------|---------|
| **Security**         | Authenticated, tamper-resistant |
| **Performance**      | Hardware-accelerated via AES-NI |
| **Integrity Check**  | Prevents tampering using a GCM tag |
| **Input size**       | Efficient for 64–128 MB data chunks |

### ✅ Why AES-GCM?

- 🔒 **AEAD**: Authenticated Encryption with Associated Data
- 🧠 **Easy to implement securely**: No need for padding
- ⚡ **Fast on most modern servers** with hardware support
- 🌐 **Widely adopted** in TLS, SSH, VPNs, etc.

#### 🔄 Compared to Other Algorithms

| Algorithm             | Speed | AEAD | Secure | Notes |
|-----------------------|-------|------|--------|-------|
| AES-CBC               | ❌     | ❌    | ⚠️      | Requires padding, vulnerable to oracle attacks |
| ChaCha20-Poly1305     | ✅     | ✅    | ✅      | Great for mobile, slower than AES-GCM on x86 |
| XChaCha20-Poly1305    | ✅     | ✅    | ✅      | Great for streaming, longer nonce |
| **AES-GCM**           | ✅✅    | ✅    | ✅      | Best performance on server CPUs (AES-NI) |

---