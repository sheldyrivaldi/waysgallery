import { Alert } from "flowbite-react";
export const LoginAlertFailed = ({ message }) => {
  return (
    <>
      <div className="flex justify-center">
        <Alert color="failure" className="flex justify-center">
          <span>
            <p>Login Failed! {message}</p>
          </span>
        </Alert>
      </div>
    </>
  );
};
