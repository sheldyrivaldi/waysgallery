import { useState } from "react";
import { useMutation } from "react-query";
import { API } from "../../config/api";

import { RegisterAlertFailed } from "../alert/RegisterAlert";

const Register = ({ showLogin, dropRegister }) => {
  const [alert, setAlert] = useState(false);
  const [alertMessage, setAlertMessage] = useState("");

  const handleShowLogin = () => {
    showLogin();
    dropRegister();
  };

  const [form, setForm] = useState({
    fullname: "",
    email: "",
    password: "",
  });

  const handleChange = (e) => {
    setForm({
      ...form,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = useMutation(async (e) => {
    try {
      e.preventDefault();

      const config = {
        headers: {
          "Content-Type": "application/json",
        },
      };

      const response = await API.post("/register", form, config);

      dropRegister();

      setForm({
        fullname: "",
        email: "",
        password: "",
      });

      setAlert(false);
    } catch (err) {
      console.log("Register Failed : ", err);
      setAlertMessage(err.response.data.message);
      setAlert(true);
    }
  });

  return (
    <>
      <div className="absolute z-[30] -top-14 w-full">{alert ? <RegisterAlertFailed message={alertMessage} /> : null}</div>
      <section id="register">
        <div className="w-[416px] h-[490px] relative z-10 py-10 px-8 rounded-lg bg-white">
          <h1 className="font-bold text-4xl text-[#2FC4B2]">Register</h1>
          <form onSubmit={(e) => handleSubmit.mutate(e)} className="w-full my-5">
            <input
              onChange={handleChange}
              type="email"
              name="email"
              placeholder="Email"
              className="text-xl my-4 text-[black] text-opacity-50 placeholder:text-[black] placeholder:opacity-50 rounded py-2 px-2 w-full bg-[#E7E7E7] border-2 border-[#5CA098]"
            />
            <input
              onChange={handleChange}
              type="password"
              name="password"
              placeholder="Password"
              className="text-xl my-4 text-[black] text-opacity-50 placeholder:text-[black] placeholder:opacity-50 rounded py-2 px-2 w-full bg-[#E7E7E7] border-2 border-[#5CA098]"
            />
            <input
              onChange={handleChange}
              type="text"
              name="fullname"
              placeholder="Full Name"
              className="text-xl mt-4  text-[black] text-opacity-50 placeholder:text-[black] placeholder:opacity-50 rounded py-2 px-2 w-full bg-[#E7E7E7] border-2 border-[#5CA098]"
            />

            <button type="submit" className="w-full rounded text-xl font-bold py-2.5 mt-8 text-white bg-[#2FC4B2] hover:bg-[#1a9b8c]">
              Register
            </button>
          </form>
          <h2 className="text-lg text-center">
            Already have an account ? Click{" "}
            <span onClick={handleShowLogin} className="font-bold cursor-pointer">
              Here
            </span>
          </h2>
        </div>
      </section>
    </>
  );
};

export default Register;
