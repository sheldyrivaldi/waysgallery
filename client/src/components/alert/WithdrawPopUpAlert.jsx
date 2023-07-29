import { Alert } from "flowbite-react";
export const WithdrawPopUpAlert = ({ message }) => {
  return (
    <>
      <div className="flex justify-center">
        <Alert color="failure" className="flex justify-center">
          <span>
            <p>{message}</p>
          </span>
        </Alert>
      </div>
    </>
  );
};
