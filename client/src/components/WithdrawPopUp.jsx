import { useQuery, useMutation } from "react-query";
import { API, setAuthToken } from "../config/api";
import { useState } from "react";

import { WithdrawPopUpAlert } from "./alert/WithdrawPopUpAlert";

const WithdrawPopUp = ({ dropWithdraw }) => {
  const [alert, setAlert] = useState(false);
  const [alertMessage, setAlertMessage] = useState("");
  const [form, setForm] = useState({
    bank_id: "",
    account_number: "",
    amount: "",
  });

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  setAuthToken(localStorage.token);

  const { data: banks } = useQuery("bankCache", async () => {
    try {
      const response = await API.get("/banks");
      return response.data.data.banks;
    } catch (err) {
      console.log("Fetching data failed!", err);
    }
  });

  const handleDropModal = () => {
    dropWithdraw();
  };

  const handleSubmitWithdraw = useMutation(async (e) => {
    e.preventDefault();
    try {
      const config = {
        headers: {
          "Content-Type": "application/json",
        },
      };
      const response = await API.post("/withdrawal", form, config);
      handleDropModal();
      setForm({
        bank_id: "",
        account_number: "",
        amount: "",
      });
    } catch (err) {
      console.log("Withdraw failed : ", err);
      setAlertMessage(err.response.data.message);
      setAlert(true);
    }
  });

  return (
    <>
      <div className="absolute z-[30] -top-14 w-full">{alert ? <WithdrawPopUpAlert message={alertMessage} /> : null}</div>
      <section id="withdraw">
        <form onSubmit={(e) => handleSubmitWithdraw.mutate(e)} className="w-[416px] h-[400px] relative z-10 py-10 mb-10 px-8 rounded-lg bg-white">
          <h1 className="font-bold text-4xl text-[#2FC4B2]">Withdraw</h1>
          <div className="w-full my-4">
            <select
              onChange={handleChange}
              name="bank_id"
              placeholder="Bank"
              className="text-xl my-4 text-[black] text-opacity-50 placeholder:text-[black] placeholder:opacity-50 rounded py-2 px-2 w-full bg-[#E7E7E7] border-2 border-[#5CA098]"
            >
              <option value="" hidden>
                Bank
              </option>
              {banks?.map((item) => {
                return <option value={item.id}>{item.name}</option>;
              })}
            </select>
            <input
              onChange={handleChange}
              type="number"
              name="account_number"
              placeholder="Account Number"
              className="text-xl my-4 text-[black] text-opacity-50 placeholder:text-[black] placeholder:opacity-50 rounded py-2 px-2 w-full bg-[#E7E7E7] border-2 border-[#5CA098]"
            />
            <input
              onChange={handleChange}
              type="number"
              name="amount"
              placeholder="Amount"
              className="text-xl my-4 text-[black] text-opacity-50 placeholder:text-[black] placeholder:opacity-50 rounded py-2 px-2 w-full bg-[#E7E7E7] border-2 border-[#5CA098]"
            />

            <button type="submit" className="w-full rounded text-xl font-bold py-2.5 mt-4 text-white bg-[#2FC4B2] hover:bg-[#1a9b8c]">
              Withdraw
            </button>
          </div>
        </form>
      </section>
    </>
  );
};

export default WithdrawPopUp;
