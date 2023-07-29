import NavbarComp from "../../components/Navbar";

import { useNavigate, useParams } from "react-router-dom";
import { useState, useEffect } from "react";
import { API, setAuthToken } from "../../config/api";
import { useQuery } from "react-query";
const Hired = () => {
  const navigate = useNavigate();
  const param = useParams();

  setAuthToken(localStorage.token);

  const { data: post } = useQuery("postHireCache", async () => {
    try {
      const response = await API.get(`/post/${param.id}`);
      return response.data.data.posts;
    } catch (err) {
      console.log("Fetching data failed!", err);
    }
  });

  const [form, setForm] = useState({
    title: "",
    description: "",
    start_date: "",
    end_date: "",
    price: "",
  });

  const handleHireBidding = async (e) => {
    e.preventDefault();
    try {
      const config = {
        headers: "application/json",
      };

      const newForm = {
        title: form.title,
        description: form.description,
        start_date: form.start_date,
        end_date: form.end_date,
        price: form.price,
        order_to_id: String(post?.user.id),
      };

      const response = await API.post("/order", newForm, config);

      const token = response?.data.data.token;
      window.snap.pay(token, {
        onSuccess: function (result) {
          navigate("/order");
        },
        onPending: function (result) {
          navigate("/order");
        },
        onError: function (result) {
          navigate("/order");
        },
        onClose: function () {
          alert("You closed the popup without finishing the payment");
        },
      });
    } catch (err) {
      console.log("Hire failed!", err);
    }
  };

  const handleNavigateHireCancel = () => {
    navigate(`/post/${param.id}`);
  };

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  useEffect(() => {
    const midtransScriptUrl = "https://app.sandbox.midtrans.com/snap/snap.js";
    const myMidtransClientKey = import.meta.env.VITE_MIDTRANS_CLIENT_KEY;

    const scriptTag = document.createElement("script");
    scriptTag.src = midtransScriptUrl;
    scriptTag.setAttribute("data-client-key", myMidtransClientKey);

    document.body.appendChild(scriptTag);
    return () => {
      document.body.removeChild(scriptTag);
    };
  }, []);

  return (
    <>
      <NavbarComp />
      <div className="px-20 pb-24 pt-36">
        <form className="w-[60%] mx-auto flex flex-col">
          <input onChange={handleChange} type="text" name="title" className="text-sm py-2 px-2 my-3 w-full rounded bg-[#E7E7E7] border-2 border-[#2FC4B2]" placeholder="Title" />
          <textarea onChange={handleChange} name="description" className="text-sm py-2 px-2 my-3 w-full h-48 rounded bg-[#E7E7E7] border-2 border-[#2FC4B2] resize-none" placeholder="Description Job"></textarea>
          <div className="flex justify-between">
            <input onChange={handleChange} type="date" name="start_date" className="text-sm py-2 px-2 my-3 w-[48%] rounded bg-[#E7E7E7] border-2 border-[#2FC4B2]" placeholder="Start Project" />
            <input onChange={handleChange} type="date" name="end_date" className="text-sm py-2 px-2 my-3 w-[48%] rounded bg-[#E7E7E7] border-2 border-[#2FC4B2]" placeholder="End Project" />
          </div>
          <input onChange={handleChange} type="number" name="price" className="text-sm py-2 px-2 my-3 w-full rounded bg-[#E7E7E7] border-2 border-[#2FC4B2] appearance-none" placeholder="Price" />

          <div className="mx-auto mt-14">
            <button onClick={handleNavigateHireCancel} type="button" className="py-2 px-8 rounded-md text-sm font-medium text-black hover:bg-[#bdbbbb] bg-[#E7E7E7]">
              Cancel
            </button>
            <button onClick={(e) => handleHireBidding(e)} type="button" className="py-2 px-8 ml-5 rounded-md text-sm font-medium hover:bg-[#1a9b8c] text-white bg-[#2FC4B2]">
              Bidding
            </button>
          </div>
        </form>
      </div>
    </>
  );
};

export default Hired;
