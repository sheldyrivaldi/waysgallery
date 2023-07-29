import { Alert } from "flowbite-react";
const ProjectAlert = () => {
  return (
    <>
      <div className="flex justify-center mt-24">
        <Alert color="success" onDismiss={() => alert("Alert dismissed!")} className="w-1/2">
          <span>
            <p>We have sent your offer, please wait for the user to accept it.</p>
          </span>
        </Alert>
      </div>
    </>
  );
};

export default ProjectAlert;
