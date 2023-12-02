import { api } from "../../../backend/api/backend";
import { APIError } from "../../../components/Error/APIError";
import { List, Title } from "@vertex-center/components";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useQuery } from "@tanstack/react-query";
import Content from "../../../components/Content/Content";
import Host from "./Host";
import CPU from "./CPU";

export default function SettingsHardware() {
    const {
        data: host,
        error: errorHost,
        isLoading: isLoadingHost,
    } = useQuery({
        queryKey: ["hardware_host"],
        queryFn: api.hardware.host,
    });

    const {
        data: cpus,
        error: errorCPUs,
        isLoading: isLoadingCPUs,
    } = useQuery({
        queryKey: ["hardware_cpus"],
        queryFn: api.hardware.cpus,
    });

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
