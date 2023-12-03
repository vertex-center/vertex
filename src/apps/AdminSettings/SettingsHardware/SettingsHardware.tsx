import { APIError } from "../../../components/Error/APIError";
import { List, Title } from "@vertex-center/components";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import Content from "../../../components/Content/Content";
import Host from "./Host";
import CPU from "./CPU";
import { useCPUs } from "../hooks/useCPUs";
import { useHost } from "../hooks/useHost";

export default function SettingsHardware() {
    const { host, errorHost, isLoadingHost } = useHost();
    const { cpus, errorCPUs, isLoadingCPUs } = useCPUs();

    const isLoading = isLoadingHost || isLoadingCPUs;
    const error = errorHost || errorCPUs;

    return (
        <Content>
            <Title variant="h2">Hardware</Title>
            <ProgressOverlay show={isLoading} />
            <APIError error={error} />

            <Title variant="h3">Host</Title>
            <Host host={host} />

            <Title variant="h3">CPUs</Title>
            <List>
                {cpus?.map((cpu, i) => (
                    <CPU key={i} cpu={cpu} />
                ))}
            </List>
        </Content>
    );
}
