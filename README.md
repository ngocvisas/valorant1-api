# Valorant Agent & Loadout API - Encore Backend

Backend API cho React app **Valorant Agent & Loadout Explorer**, build bằng [Encore](https://encore.dev).

---

 Features
- **GET /agents** – Lấy danh sách agents (filter/search)
- **GET /weapons** – Lấy danh sách weapons (filter/maxCost/search)
- **GET /health** – Health check
- **GET /stats** – Thống kê tổng quan
- **POST /loadouts** – (auth) Tạo loadout
- **GET /loadouts** – (auth) Lấy loadouts user

---

 Setup Local
```bash
# Clone project
git clone https://github.com/<your-username>/valorant1-api.git
cd valorant1-api

# Run Encore dev
encore run

