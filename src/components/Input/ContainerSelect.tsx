import Select, { SelectOption, SelectValue } from "./Select";
import { Container, ContainerQuery } from "../../models/container";
import { api } from "../../backend/api/backend";
import Progress from "../Progress";
import ServiceLogo from "../ServiceLogo/ServiceLogo";
import { useQuery } from "@tanstack/react-query";
import { APIError } from "../Error/APIError";

type Props = {
    container?: Container;
    onChange?: (container?: Container) => void;

    query?: ContainerQuery;
};

export default function ContainerSelect(props: Readonly<Props>) {
    const { container, onChange, query } = props;

    const queryContainers = useQuery({
        queryKey: ["containers", query],
        queryFn: () => api.vxContainers.containers.search(query),
    });
    const { data: containers, isLoading, error } = queryContainers;

    if (isLoading) {
        return <Progress infinite />;
    }

    if (error) {
        return <APIError error={error} />;
    }

    const onContainerChange = (uuid: any) => {
        const container = containers?.[uuid];
        onChange?.(container);
    };

    const value = (
        <SelectValue>
            {container && <ServiceLogo service={container?.service} />}
            {container?.display_name ??
                container?.service?.name ??
                "Select an container"}
        </SelectValue>
    );

    return (
        // @ts-ignore
        <Select onChange={onContainerChange} value={value}>
            <SelectOption value="">None</SelectOption>
            {Object.entries(containers ?? [])?.map(([, container]) => (
                <SelectOption key={container.uuid} value={container.uuid}>
                    <ServiceLogo service={container?.service} />
                    {container?.display_name ?? container?.service?.name}
                </SelectOption>
            ))}
        </Select>
    );
}
