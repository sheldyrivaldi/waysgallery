export function FormatToRupiah(number) {
  // Menggunakan metode toLocaleString() untuk memformat angka menjadi mata uang Rupiah
  let formattedNumber = "Rp 0";
  if (number != 0) {
    formattedNumber = new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR" }).format(number);
  }

  // Menghilangkan ".00" pada bagian desimal jika ada
  formattedNumber = formattedNumber.replace(/\,00$/, "");

  return formattedNumber;
}
