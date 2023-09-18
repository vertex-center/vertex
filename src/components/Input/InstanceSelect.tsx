import Select, { SelectOption, SelectValue } from "./Select";
import { Instance, InstanceQuery, Instances } from "../../models/instance";
import { useFetch } from "../../hooks/useFetch";
import { api } from "../../backend/backend";
import Progress from "../Progress";
import ServiceLogo from "../ServiceLogo/ServiceLogo";

type Props = {
    instance?: Instance;
    onChange?: (instance?: Instance) => void;

    query?: InstanceQuery;
};

export default function InstanceSelect(props: Props) {
    const { instance, onChange, query } = props;

    const search = () => api.instances.search(query).catch(console.error);
    const { data: instances, loading } = useFetch<Instances>(search);

    if (loading) {
        return <Progress infinite />;
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
