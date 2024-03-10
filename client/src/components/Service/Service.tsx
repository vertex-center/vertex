import { Template as TemplateModel } from "../../apps/Containers/backend/template";
import ServiceLogo from "../ServiceLogo/ServiceLogo";
import { Card, Title, Vertical } from "@vertex-center/components";

type Props = {
    template: TemplateModel;
    onInstall: () => void;
};

export default function Service(props: Readonly<Props>) {
    const { template, onInstall } = props;

    return (
        <Card onClick={onInstall}>
            <Vertical gap={40}>
                <ServiceLogo icon={template?.icon} color={template?.color} />
                <Title>{template?.name}</Title>
            </Vertical>
        </Card>
    );
}
