import Hardware from "../../../components/Hardware/Hardware";
import { api } from "../../../backend/api/backend";
import { APIError } from "../../../components/Error/APIError";
import { List, Title } from "@vertex-center/components";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useQuery } from "@tanstack/react-query";
import Content from "../../../components/Content/Content";

export default function SettingsHardware() {
    const {
        data: hardware,
        error,
        isLoading,
    } = useQuery({
        queryKey: ["hardware"],
        queryFn: api.hardware,
    });

    return (
        <Content>
            <Title variant="h2">Hardware</Title>
            <ProgressOverlay show={isLoading} />
            <APIError error={error} />
            <List>
                <Hardware hardware={hardware} />
            </List>
        </Content>
    );
}
