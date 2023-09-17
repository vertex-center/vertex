import Select, { Option } from "./Select";
import { Instance, InstanceQuery } from "../../models/instance";
import { useFetch } from "../../hooks/useFetch";
import { api } from "../../backend/backend";
import Progress from "../Progress";

type Props = {
    instance?: Instance;
    onChange?: (instance?: Instance) => void;

    query?: InstanceQuery;
};

export default function InstanceSelect(props: Props) {
    const { instance, onChange, query } = props;

    const search = () => api.instances.search(query).catch(console.error);
    const { data: instances, loading } = useFetch(search);

    if (loading) {
        return <Progress infinite />;
    }

    const onInstanceChange = (e: any) => {
        const uuid = e.target.value;
        const instance = instances?.[uuid];
        console.log(instance);
        onChange?.(instance);
    };

    return (
        <Select onChange={onInstanceChange} value={instance?.uuid}>
            <Option value=""></Option>
            {Object.entries(instances ?? [])?.map(([, instance]) => (
                <Option key={instance.uuid} value={instance.uuid}>
                    {instance.name}
                </Option>
            ))}
        </Select>
    );
}
