# Login Handling Components

Berikut adalah komponen-komponen utama yang menangani login di LobeChat:

## 1. Komponen UI Utama
- **[AuthSignInBox.tsx](./src/app/[variants]/(auth)/next-auth/signin/AuthSignInBox.tsx)**
  Komponen ini adalah UI utama untuk sign-in menggunakan **Next-Auth**. Ini menampilkan daftar provider (seperti Google, GitHub, dll) dan menangani klik tombol sign-in.
- **[UserLoginOrSignup/Community.tsx](./src/features/User/UserLoginOrSignup/Community.tsx)**
  Komponen ini menyediakan tombol "Login atau Daftar" (Login or Signup) yang biasanya muncul di sidebar atau menu profil.
- **[ClerkLogin/index.tsx](./src/features/Conversation/Error/ClerkLogin/index.tsx)**
  Komponen ini muncul saat terjadi error terkait autentikasi (khusus untuk **Clerk**) di dalam percakapan, memberikan tombol bagi pengguna untuk masuk.

## 2. Logika Utama (Store Action)
- **[user/slices/auth/action.ts](./src/store/user/slices/auth/action.ts)**
  Fungsi `openLogin` di file ini adalah logika sentral yang dipanggil oleh berbagai komponen UI di atas. Fungsi ini menentukan apakah harus menggunakan Clerk atau Next-Auth berdasarkan konfigurasi environment (`enableClerk` atau `enableNextAuth`).

### Ringkasan Alur:
1. Pengguna mengklik tombol login di `UserLoginOrSignup`.
2. `UserLoginOrSignup` memanggil `openLogin` dari `useUserStore`.
3. `openLogin` memicu alur autentikasi (baik redirect ke halaman Clerk atau membuka halaman sign-in Next-Auth).
4. Jika menggunakan Next-Auth, pengguna akan diarahkan ke halaman yang menggunakan `AuthSignInBox`.
