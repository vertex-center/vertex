import Select, { SelectOption, SelectValue } from "./Select";
import { Instance, InstanceQuery } from "../../models/instance";
import { api } from "../../backend/api/backend";
import Progress from "../Progress";
import ServiceLogo from "../ServiceLogo/ServiceLogo";
import { useQuery } from "@tanstack/react-query";
import { APIError } from "../Error/APIError";

type Props = {
    instance?: Instance;
    onChange?: (instance?: Instance) => void;

    query?: InstanceQuery;
};

export default function InstanceSelect(props: Readonly<Props>) {
    const { instance, onChange, query } = props;

    const queryInstances = useQuery({
        queryKey: ["instances", query],
        queryFn: () => api.vxInstances.instances.search(query),
    });
    const { data: instances, isLoading, error } = queryInstances;

    if (isLoading) {
        return <Progress infinite />;
    }

    if (error) {
        return <APIError error={error} />;
    }

    const onInstanceChange = (uuid: any) => {
        const instance = instances?.[uuid];
        onChange?.(instance);
    };

    const value = (
        <SelectValue>
            {instance && <ServiceLogo service={instance?.service} />}
            {instance?.display_name ??
                instance?.service?.name ??
                "Select an instance"}
        </SelectValue>
    );

    return (
        // @ts-ignore
        <Select onChange={onInstanceChange} value={value}>
            <SelectOption value="">None</SelectOption>
            {Object.entries(instances ?? [])?.map(([, instance]) => (
                <SelectOption key={instance.uuid} value={instance.uuid}>
                    <ServiceLogo service={instance?.service} />
                    {instance?.display_name ?? instance?.service?.name}
                </SelectOption>
            ))}
        </Select>
    );
}
