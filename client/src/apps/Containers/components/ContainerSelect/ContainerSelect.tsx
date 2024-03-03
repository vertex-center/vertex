import { SelectField, SelectOption } from "@vertex-center/components";
import Progress from "../../../../components/Progress";
import ServiceLogo from "../../../../components/ServiceLogo/ServiceLogo";
import { useQuery } from "@tanstack/react-query";
import { APIError } from "../../../../components/Error/APIError";
import { Fragment } from "react";
import { API } from "../../backend/api";
import { Container, ContainerFilters } from "../../backend/models";

type Props = {
    container?: Container;
    onChange?: (container?: Container) => void;
    filters?: ContainerFilters;
};

export default function ContainerSelect(props: Readonly<Props>) {
    const { container, onChange, filters } = props;

    const queryContainers = useQuery({
        queryKey: ["containers", filters],
        queryFn: () => API.getContainers(filters),
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
        <Fragment>
            {container && <ServiceLogo service={container?.service} />}
            {container?.name ?? "Select an container"}
        </Fragment>
    );

    return (
        // @ts-ignore
        <SelectField onChange={onContainerChange} value={value}>
            <SelectOption value="">None</SelectOption>
            {Object.entries(containers ?? [])?.map(([, c]) => (
                <SelectOption key={c?.id} value={c?.id}>
                    <ServiceLogo service={c?.service} />
                    {c?.name}
                </SelectOption>
            ))}
        </SelectField>
    );
}
