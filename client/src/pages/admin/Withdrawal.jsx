import { TableHeadCell } from "flowbite-react/lib/esm/components/Table/TableHeadCell";
import NavbarComp from "../../components/Navbar";
import { Table } from "flowbite-react";
import { FormatToRupiah } from "../../utils/FormatRupiah";
import { useQuery } from "react-query";
import { API, setAuthToken } from "../../config/api";
import { useState, useEffect } from "react";
import Accept from "../../assets/accept.svg";

const Withdrawal = () => {
  const [isSubmit, setIsSubmit] = useState(false);

  setAuthToken(localStorage.token);

  const { data: withdrawals, refetch } = useQuery("withdrawalCache", async () => {
    try {
      const response = await API.get("/withdrawals");
      return response.data.data.withdrawals;
    } catch (err) {
      console.log("Failed fetching data!", err);
    }
  });

  const handleAccept = async (e, id) => {
    e.preventDefault();
    try {
      const config = {
        headers: {
          "Content-Type": "application/json",
        },
      };
      const response = await API.patch(
        `/withdrawal/${id}`,
        {
          status: "Success",
        },
        config
      );

      setIsSubmit(true);
    } catch (err) {
      console.log("Failed accept withdrawal!", err);
    }
  };
  useEffect(() => {
    refetch();
    setIsSubmit(false);
  }, [isSubmit]);

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

  const handleAction = (status, id) => {
    switch (status) {
      case "Pending":
        return (
          <Table.Cell className="flex justify-center">
            <button onClick={(e) => handleAccept(e, id)} className="text-sm font-bold px-3 py-1 rounded text-white bg-[#00D1FF] hover:bg-[#07afd4]" type="button">
              Accept
            </button>
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
            <button className="text-sm font-bold px-3 py-1 rounded text-white bg-[#00D1FF] hover:bg-[#07afd4]" type="button">
              Accept
            </button>
          </Table.Cell>
        );
    }
  };

  useEffect(() => {
    refetch();
    setIsSubmit(false);
  }, [isSubmit]);
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
                      {handleAction(item.status, item.id)}
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

export default Withdrawal;
