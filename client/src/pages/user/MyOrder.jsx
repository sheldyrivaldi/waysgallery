import { TableHeadCell } from "flowbite-react/lib/esm/components/Table/TableHeadCell";
import NavbarComp from "../../components/Navbar";
import { Table } from "flowbite-react";
import Accept from "../../assets/accept.svg";
import Cancel from "../../assets/cancel.svg";
import WaitingAccept from "../../assets/waiting-accept.svg";
import ProjectPopUp from "../../components/ProjectPopUp";
import WithdrawPopUp from "../../components/WithdrawPopUp";

import { useNavigate } from "react-router-dom";
import { useState, useEffect, useContext } from "react";
import { useQuery, useMutation } from "react-query";
import { API, setAuthToken } from "../../config/api";
import { UserContext } from "../../context/Context";
import { FormatToRupiah } from "../../utils/FormatRupiah";
import { Modal } from "flowbite-react";

const MyOrder = () => {
  const [select, setSelect] = useState("my-order");
  const [state] = useContext(UserContext);
  const [openModal, setOpenModal] = useState();
  const props = { openModal, setOpenModal };
  const [openModalWithdraw, setOpenModalWithdraw] = useState();
  const propsWithdraw = { openModalWithdraw, setOpenModalWithdraw };
  const [order, setOrder] = useState({});
  const navigate = useNavigate();

  const handlerNavigateMyOffer = () => {
    if (select == "my-offer") {
      navigate("/offer");
    }
  };

  setAuthToken(localStorage.token);

  const { data: orders } = useQuery("orderCache", async () => {
    const response = await API.get(`/orders/client/${state.user.id}`);
    return response.data.data.orders;
  });

  const { data: user } = useQuery("userOrderCache", async () => {
    const response = await API.get(`/user/${state.user.id}`);
    return response.data.data.users;
  });

  useEffect(() => {
    handlerNavigateMyOffer();
  }, [select]);

  const handleNavigateViewProject = (id) => {
    navigate(`/order/${id}`);
  };

  const handleOpenModal = (order) => {
    props.setOpenModal("dismissible");
    setOrder(order);
  };

  const handleOpenModalWithdraw = () => {
    propsWithdraw.setOpenModalWithdraw("dismissible");
  };

  const handleStatus = (status) => {
    switch (status) {
      case "Waiting Accept":
        return <Table.Cell className="text-[#FF9900]">Waiting Accept</Table.Cell>;
      case "Success":
        return <Table.Cell className="text-[#78A85A]">Success</Table.Cell>;
      case "Cancel":
        return <Table.Cell className="text-[#E83939]">Cancel</Table.Cell>;
      case "Project is Complete":
        return <Table.Cell className="text-[#00D1FF]">Project is Complete</Table.Cell>;
      case "On Progress":
        return <Table.Cell className="text-[#00D1FF]">On Progress</Table.Cell>;
      case "Revision":
        return <Table.Cell className="text-[#FF9900]">Revision</Table.Cell>;
      default:
        return <Table.Cell className="text-[#FF9900]">Waiting Accept</Table.Cell>;
    }
  };

  const handleAction = (status, id) => {
    switch (status) {
      case "Waiting Accept":
        return (
          <Table.Cell className="flex justify-center">
            <img src={WaitingAccept} alt="waiting-accept" />
          </Table.Cell>
        );
      case "Success":
        return (
          <Table.Cell className="flex justify-center">
            <button onClick={() => handleNavigateViewProject(id)} type="button" className="text-sm px-3 py-1 font-bold rounded text-white bg-[#0ACF83] hover:bg-[#09bd78]">
              View Project
            </button>
          </Table.Cell>
        );
      case "Cancel":
        return (
          <Table.Cell className="flex justify-center">
            <img src={Cancel} alt="cancel" />
          </Table.Cell>
        );
      case "Project is Complete":
        return (
          <Table.Cell className="flex justify-center">
            <button onClick={() => handleNavigateViewProject(id)} type="button" className="text-sm px-3 py-1 font-bold rounded text-white bg-[#0ACF83] hover:bg-[#09bd78]">
              View Project
            </button>
          </Table.Cell>
        );
      case "On Progress":
        return <Table.Cell className="text-[#00D1FF]">On Progress</Table.Cell>;
      case "Revision":
        return (
          <Table.Cell className="flex justify-center">
            <img src={WaitingAccept} alt="revision" />
          </Table.Cell>
        );
      default:
        return <Table.Cell className="text-[#FF9900]">Waiting Accept</Table.Cell>;
    }
  };

  return (
    <>
      <NavbarComp />
      <div className="px-20 py-24">
        <div className="w-full h-14 flex justify-between items-start">
          <select
            onChange={(e) => setSelect(e.target.value)}
            class="bg-gray-50 border mt-5 border-gray-300 text-gray-900 text-xs rounded focus:ring-blue-500 focus:border-blue-500 block w-28 px-3 py-2 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
          >
            <option value="my-order">My Order</option>
            <option value="my-offer">My Offer</option>
          </select>

          <div className="flex flex-col justify-center items-end">
            <h1 className="text-lg font-bold">
              Balance : <span className="text-[#78A85A]">{FormatToRupiah(user?.balance)}</span>
            </h1>
            <button onClick={() => handleOpenModalWithdraw()} className="text-sm mt-5 px-3 py-1 font-bold rounded text-white bg-[#0ACF83] hover:bg-[#09bd78]" type="button">
              Withdraw
            </button>
          </div>
        </div>
        <div className="mt-20 px-5">
          <Table>
            <Table.Head>
              <TableHeadCell className="bg-[#E5E5E5]">No</TableHeadCell>
              <TableHeadCell className="bg-[#E5E5E5]">Vendor</TableHeadCell>
              <TableHeadCell className="bg-[#E5E5E5]">Order</TableHeadCell>
              <TableHeadCell className="bg-[#E5E5E5]">Start Project</TableHeadCell>
              <TableHeadCell className="bg-[#E5E5E5]">End Project</TableHeadCell>
              <TableHeadCell className="bg-[#E5E5E5]">Status</TableHeadCell>
              <TableHeadCell className="text-center bg-[#E5E5E5]">Action</TableHeadCell>
            </Table.Head>
            <Table.Body>
              {orders && orders?.length > 0 ? (
                orders?.map((item) => (
                  <Table.Row>
                    <Table.Cell className="text-center">{item.id}</Table.Cell>
                    <Table.Cell>{item.order_to.fullname}</Table.Cell>
                    <Table.Cell>
                      <span onClick={() => handleOpenModal(item)} className="cursor-pointer text-[#0D33B9]">
                        {item.title}
                      </span>
                    </Table.Cell>
                    <Table.Cell>{item.start_date}</Table.Cell>
                    <Table.Cell>{item.end_date}</Table.Cell>
                    {handleStatus(item.status)}
                    {handleAction(item.status, item.id)}
                  </Table.Row>
                ))
              ) : (
                <div className="text-lg relative left-[30rem] font-bold text-center my-10">No Order</div>
              )}
            </Table.Body>
          </Table>
        </div>
      </div>
      <Modal dismissible show={props.openModal === "dismissible"} onClose={() => props.setOpenModal(undefined)} size="xl">
        <ProjectPopUp order={order} />
      </Modal>
      <Modal dismissible show={propsWithdraw.openModalWithdraw === "dismissible"} onClose={() => propsWithdraw.setOpenModalWithdraw(undefined)} size="md">
        <WithdrawPopUp dropWithdraw={() => propsWithdraw.setOpenModalWithdraw(undefined)} />
      </Modal>
    </>
  );
};

export default MyOrder;
