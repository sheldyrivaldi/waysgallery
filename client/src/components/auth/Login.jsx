import { useMutation } from "react-query";
import { useContext, useState } from "react";
import { useNavigate } from "react-router-dom";
import { UserContext } from "../../context/Context";
import { API, setAuthToken } from "../../config/api";

import { LoginAlertFailed } from "../alert/LoginAlert";

const Login = ({ showRegister, dropLogin }) => {
  const [alert, setAlert] = useState(false);
  const [alertMessage, setAlertMessage] = useState("");
  const handleShowRegister = () => {
    showRegister();
    dropLogin();
  };

  const [_, dispatch] = useContext(UserContext);
  const navigate = useNavigate();
  const [form, setForm] = useState({
    email: "",
    password: "",
  });

  const config = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = useMutation(async (e) => {
    try {
      e.preventDefault();

      const response = await API.post("/login", form, config);
      dispatch({
        type: "LOGIN_SUCCESS",
        payload: response.data.data.user,
      });

      setAuthToken(localStorage.token);

      if (response.data.data.user?.role === "admin") {
        navigate("/admin");
      } else {
        navigate("/home");
      }
      dropLogin();

      setForm({
        email: "",
        password: "",
      });
    } catch (err) {
      console.log("login failed : ", err);
      setAlertMessage(err.response.data.message);
      setAlert(true);
    }
  });

  return (
    <>
      <div className="absolute z-[30] -top-14 w-full">{alert ? <LoginAlertFailed message={alertMessage} /> : null}</div>
      <section id="login">
        <div className="w-[416px] h-[400px] relative z-10 py-10 px-8 rounded-lg bg-white">
          <h1 className="font-bold text-4xl text-[#2FC4B2]">Login</h1>
          <form onSubmit={(e) => handleSubmit.mutate(e)} className="w-full my-4">
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
            <button type="submit" className="w-full rounded text-xl font-bold py-2.5 mt-4 text-white bg-[#2FC4B2] hover:bg-[#1a9b8c]">
              Login
            </button>
          </form>
          <h2 className="text-lg text-center">
            Don't have an account ? Click{" "}
            <span onClick={handleShowRegister} className="font-bold cursor-pointer">
              Here
            </span>
          </h2>
        </div>
      </section>
    </>
  );
};

export default Login;
