import { TableHeadCell } from "flowbite-react/lib/esm/components/Table/TableHeadCell";
import NavbarComp from "../../components/Navbar";
import { Table } from "flowbite-react";
import { FormatToRupiah } from "../../utils/FormatRupiah";
import { useQuery } from "react-query";
import { API, setAuthToken } from "../../config/api";
import { useState, useEffect, useContext } from "react";
import { UserContext } from "../../context/Context";
import Accept from "../../assets/accept.svg";
import Pending from "../../assets/waiting-accept.svg";

const WithdrawalUser = () => {
  const [isSubmit, setIsSubmit] = useState(false);
  const [state] = useContext(UserContext);

  setAuthToken(localStorage.token);

  const { data: withdrawals } = useQuery("withdrawalUserCache", async () => {
    try {
      const response = await API.get(`/withdrawals/${state.user.id}`);
      return response.data.data.withdrawals;
    } catch (err) {
      console.log("Failed fetching data!", err);
    }
  });

  const handleStatus = (status) => {
    switch (status) {
      case "Pending":
        return <Table.Cell className="text-[#FF9900]">Pending</Table.Cell>;
      case "Success":
        return <Table.Cell className="text-[#78A85A]">Success</Table.Cell>;
      default:
        return <Table.Cell className="text-[#FF9900]">Pending</Table.Cell>;
    }
  };

  const handleAction = (status) => {
    switch (status) {
      case "Pending":
        return (
          <Table.Cell className="flex justify-center">
            <img src={Pending} alt="pending" />
          </Table.Cell>
        );
      case "Success":
        return (
          <Table.Cell className="flex justify-center">
            <img src={Accept} alt="success" />
          </Table.Cell>
        );
      default:
        return (
          <Table.Cell className="flex justify-center">
            <img src={Pending} alt="pending" />
          </Table.Cell>
        );
    }
  };
  return (
    <>
      <NavbarComp />
      <div className="px-20 py-24">
        <div className="mt-20 px-5">
          <Table>
            <Table.Head>
              <TableHeadCell className="bg-[#E5E5E5]">No</TableHeadCell>
              <TableHeadCell className="bg-[#E5E5E5]">User</TableHeadCell>
              <TableHeadCell className="bg-[#E5E5E5]">Amount</TableHeadCell>
              <TableHeadCell className="bg-[#E5E5E5]">Bank</TableHeadCell>
              <TableHeadCell className="bg-[#E5E5E5]">Account Number</TableHeadCell>
              <TableHeadCell className="bg-[#E5E5E5]">Status</TableHeadCell>
              <TableHeadCell className="bg-[#E5E5E5] text-center">Action</TableHeadCell>
            </Table.Head>
            <Table.Body>
              {withdrawals && withdrawals?.length > 0 ? (
                withdrawals?.map((item) => {
                  return (
                    <Table.Row>
                      <Table.Cell>{item.id}</Table.Cell>
                      <Table.Cell>{item.user.fullname}</Table.Cell>
                      <Table.Cell>{FormatToRupiah(item.amount)}</Table.Cell>
                      <Table.Cell>{item.bank.name}</Table.Cell>
                      <Table.Cell>{item.account_number}</Table.Cell>
                      {handleStatus(item.status)}
                      {handleAction(item.status)}
                    </Table.Row>
                  );
                })
              ) : (
                <div className="text-lg relative left-[30rem] font-bold text-center my-10">No Data</div>
              )}
            </Table.Body>
          </Table>
        </div>
      </div>
    </>
  );
};

export default WithdrawalUser;
