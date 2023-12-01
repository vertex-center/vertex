import Content from "../../../../components/Content/Content";
import { Title } from "@vertex-center/components";
import useUser from "../../hooks/useUser";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";

export default function AccountInfo() {
    const { user, isLoadingUser, errorUser } = useUser();
    return (
        <Content>
            <Title variant="h2">Information</Title>
            <ProgressOverlay show={isLoadingUser} />
            <APIError error={errorUser} />
            <div>{user?.toString()}</div>
        </Content>
    );
}
