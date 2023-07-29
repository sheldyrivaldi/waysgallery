export default function formatDate(date) {
  const [day, month, year] = date.split("-");

  const newDate = new Date(year, month - 1, day);

  const daysOfWeek = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
  const dayOfWeek = daysOfWeek[newDate.getDay()];

  const monthNames = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];
  const monthName = monthNames[month - 1];

  return [dayOfWeek, day, monthName, year];
}
