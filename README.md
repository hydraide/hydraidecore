![HydrAIDE Core](images/hydraide-banner.jpg)

# 🧠 HydrAIDE Core – The Adaptive, Intelligent Data Engine

![Go](https://img.shields.io/badge/built%20with-Go-00ADD8?style=for-the-badge&logo=go)
![Status](https://img.shields.io/badge/status-Production%20Ready-brightgreen?style=for-the-badge)
![Version](https://img.shields.io/badge/version-2.0-informational?style=for-the-badge)
![Speed](https://img.shields.io/badge/Access-O(1)%20Always-ff69b4?style=for-the-badge)
![License: Core](https://img.shields.io/badge/license-SSPL--1.0--Custom-red?style=for-the-badge)

> **HydrAIDE is a zero-waste, real-time data engine. Like a database, but built for reactive systems**

---

## ⚙️ Do I need to build the Core myself?

**Nope. You don’t have to build anything manually.**

The HydrAIDE Core is already included in the official HydrAIDE server,  
which we distribute as a **ready-to-run Docker container**.

This repository is intended for:

- developers who want to **debug or extend** the engine,
- teams embedding HydrAIDE Core **as an in-process database** in their own Go apps,
- or advanced users building for **IoT or edge** deployments.

---

## 🔧 Why is this repo separate?

Because the Core can work **without any server at all.**  
It’s a **zero-daemon**, zero-overhead data engine, which can be compiled straight into your Go application.

That means:

- No background processes.
- No runtime listeners.
- No unnecessary resource usage.

If your use case requires an **in-process embedded database**,  
you can import the Core and use it **just like a package** — no fuss.

But if you want to run HydrAIDE as a **networked external engine**,  
just use our **official Docker container**, which includes the full Core engine,  
plus the gRPC interface.

→ [https://github.com/hydraide/hydraide](https://github.com/hydraide/hydraide)

