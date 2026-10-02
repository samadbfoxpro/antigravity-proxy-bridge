# 🛰️ Antigravity Proxy Bridge

<p align=center>
  <b>A lightweight, standalone CLI proxy manager & connectivity tester for Google Antigravity IDE & Antigravity 2.0.</b>
  <br>
  <i>Connect Antigravity to local VPN/Proxy clients (v2rayN, NekoBox, Clash, Hiddify) without tunneling your entire OS.</i>
</p>

<p align=center>
  <img src=https://img.shields.io/badge/Platform-Windows-blue?logo=windows alt=Windows>
  <img src=https://img.shields.io/badge/Built%20With-Go-00ADD8?logo=go alt=Go>
  <img src=https://img.shields.io/badge/Zero%20Dependency-Standalone-brightgreen alt=Standalone>
  <img src=https://img.shields.io/badge/License-MIT-yellow.svg alt=License>
</p>

---

## 📖 معرفی به فارسی (Overview)

هنگام کار با **Google Antigravity** در مناطق با دسترسی محدود، معمولاً مجبور به استفاده از حالت VPN Mode / TUN برای کل سیستم‌عامل هستیم که سرعت سایر برنامه‌ها و دانلودها را مختل می‌کند. 

**Antigravity Proxy Bridge** یک ابزار سبک و تک‌فایلی (.exe) بدون نیاز به نصب هیچ پیش‌نیازی است که ترافیک Antigravity (هم نسخه IDE و هم نسخه 2.0) را مستقیماً به پورت لوکال کلاینت فیلترشکن شما هدایت می‌کند.

---

## ✨ ویژگی‌های کلیدی (Key Features)

- ⚡ **تک‌فایلی و پرتابل (Zero Dependency):** کامپایل شده به زبان Go؛ بدون نیاز به نصب پایتون یا هر کتابخانه دیگر.
- 🟢 **نمایش زنده وضعیت (Real-time Status):** تشخیص فوری وضعیت پروکسی با رنگ‌های واضح در منو.
- 🎯 **تفکیک‌پذیری کامل:**
  - فعال‌سازی فقط برای **Antigravity IDE**
  - فعال‌سازی فقط برای **Antigravity 2.0**
  - فعال‌سازی همزمان برای **هر دو**
- 📶 **تست پینگ و لتنسی (Latency Tester):** بررسی اتصال زنده به سرورهای گوگل و اندازه‌گیری زمان پاسخ بر حسب میلی‌ثانیه.
- ⚙️ **تغییر آسان تنظیمات:** امکان تغییر سریع IP، پورت و پروتکل (HTTP / SOCKS5) بدون نیاز به دستکاری کد.
- 🛡️ **بازگشت به حالت اول (One-Click Reset):** حذف کامل و بدون ردپای تنظیمات پروکسی تنها با یک گزینه.

---

## 🔌 پورت‌های پیش‌فرض کلاینت‌ها (Common Client Ports)

| کلاینت فیلترشکن | پروتکل | پورت محلی |
| :--- | :--- | :--- |
| **v2rayN** | HTTP / SOCKS5 | 10809 (HTTP) / 10808 (SOCKS) |
| **NekoRay / NekoBox** | HTTP / SOCKS5 | 2080 (HTTP) / 2081 (SOCKS) |
| **Hiddify** | HTTP | 2334 |
| **Clash / Sing-box** | Mixed HTTP/SOCKS | 7890 |

---

## 🚀 نحوه استفاده (Quick Start)

1. مطمئن شوید برنامه فیلترشکن شما باز و فعال است (نیازی به روشن بودن حالت TUN یا VPN سیستم نیست).
2. فایل AntigravityProxyManager.exe را اجرا کنید.
3. در منوی برنامه:
   - کلید 1 را برای اعمال روی **IDE** بزنید.
   - کلید 2 را برای اعمال روی **Antigravity 2.0** بزنید.
   - یا کلید 3 را برای اعمال روی **هر دو** بزنید.
4. با زدن کلید 5 می‌توانید از برقراری اتصال و سلامت پروکسی اطمینان حاصل کنید.
5. نرم‌افزار Antigravity را یک‌بار ببندید و دوباره باز کنید.

---

## 🛠️ ساخت از سورس (Build from Source)

اگر می‌خواهید خودتان برنامه را کامپایل کنید:

`ash
git clone https://github.com/your-username/antigravity-proxy-bridge.git
cd antigravity-proxy-bridge
go build -ldflags=-s -w -o AntigravityProxyManager.exe main.go
`

---

## 📄 License
This project is open-source and available under the [MIT License](LICENSE).
