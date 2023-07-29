import { FormatToRupiah } from "../utils/FormatRupiah";

const ProjectPopUp = ({ order }) => {
  return (
    <>
      <div className="w-full h-full">
        <div className="flex flex-col justify-between w-full h-56 px-5 py-5 rounded bg-white">
          <h2 className="text-black text-opacity-50 text-xl">Title : {order.title}</h2>
          <h2 className="text-black text-opacity-50 text-xl">Description : {order.description}</h2>
          <h2 className="font-bold text-xl text-[#00E016]">Price : {FormatToRupiah(order.price)}</h2>
        </div>
      </div>
    </>
  );
};

export default ProjectPopUp;
