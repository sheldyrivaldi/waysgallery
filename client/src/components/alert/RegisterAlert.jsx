import { Alert } from "flowbite-react";
export const RegisterAlertFailed = ({ message }) => {
  return (
    <>
      <div className="flex justify-center">
        <Alert color="failure" className="flex justify-center">
          <span>
            <p>Register Failed! {message}</p>
          </span>
        </Alert>
      </div>
    </>
  );
};
