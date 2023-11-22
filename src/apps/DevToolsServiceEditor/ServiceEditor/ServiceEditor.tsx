import styles from "./ServiceEditor.module.sass";
import { Horizontal, Vertical } from "../../../components/Layouts/Layouts";
import NoItems from "../../../components/NoItems/NoItems";
import {
    Button,
    Code,
    List,
    ListItem,
    MaterialIcon,
    SelectField,
    SelectOption,
    TextField,
    Title,
    useTitle,
} from "@vertex-center/components";
import classNames from "classnames";
import Spacer from "../../../components/Spacer/Spacer";
import {
    Controller,
    SubmitHandler,
    useFieldArray,
    useForm,
} from "react-hook-form";
import * as yup from "yup";
import { object } from "yup";
import { yupResolver } from "@hookform/resolvers/yup";
import { Fragment, useState } from "react";
import Card from "../../../components/Card/Card";
import { api } from "../../../backend/api/backend";
import { produce } from "immer";
import Content from "../../../components/Content/Content";

function EnvironmentInputs({ control, register, errors, i }) {
    return (
        <Fragment>
            <Controller
                render={({ field }) => {
                    const value = environmentTypes.find(
                        (e) => e.value === field.value
                    )?.label;
                    return (
                        <SelectField
                            id={`environment.${i}.type`}
                            label="Type"
                            {...register(`environment.${i}.type`)}
                            // @ts-ignore
                            value={value}
                            onChange={field.onChange}
                            required
                        >
                            {environmentTypes.map((e) => (
                                <SelectOption key={e.value} value={e.value}>
                                    {e.label}
                                </SelectOption>
                            ))}
                        </SelectField>
                    );
                }}
                name={`environment.${i}.type`}
                control={control}
            />
            <TextField
                id={`environment.${i}.name`}
                label="Name"
                {...register(`environment.${i}.name`)}
                aria-invalid={errors?.name ? "true" : "false"}
                error={errors?.name?.message}
                required
            />
            <TextField
                id={`environment.${i}.display_name`}
                label="Display name"
                {...register(`environment.${i}.display_name`)}
                aria-invalid={errors?.display_name ? "true" : "false"}
                error={errors?.display_name?.message}
                required
            />
            <TextField
                id={`environment.${i}.default`}
                label="Default value"
                {...register(`environment.${i}.default`)}
                aria-invalid={errors?.default ? "true" : "false"}
                error={errors?.default?.message}
            />
            <TextField
                id={`environment.${i}.description`}
                label="Description"
                {...register(`environment.${i}.description`)}
                aria-invalid={errors?.description ? "true" : "false"}
                error={errors?.description?.message}
            />
        </Fragment>
    );
}

function UrlInputs({ control, register, errors, i }) {
    return (
        <Fragment>
            <TextField
                id={`urls.${i}.name`}
                label="Name"
                {...register(`urls.${i}.name`)}
                aria-invalid={errors?.name ? "true" : "false"}
                error={errors?.name?.message}
                required
            />
            <TextField
                id={`urls.${i}.port`}
                label="Port"
                {...register(`urls.${i}.port`)}
                aria-invalid={errors?.port ? "true" : "false"}
                error={errors?.port?.message}
                required
            />
            <Controller
                render={({ field }) => {
                    const value = urlKinds.find(
                        (e) => e.value === field.value
                    )?.label;
                    return (
                        <SelectField
                            id={`urls.${i}.kind`}
                            label="Kind"
                            {...register(`urls.${i}.kind`)}
                            // @ts-ignore
                            value={value}
                            onChange={field.onChange}
                            required
                        >
                            {urlKinds.map((e) => (
                                <SelectOption key={e.value} value={e.value}>
                                    {e.label}
                                </SelectOption>
                            ))}
                        </SelectField>
                    );
                }}
                name={`urls.${i}.kind`}
                control={control}
            />
            <TextField
                id={`urls.${i}.ping`}
                label="Ping"
                {...register(`urls.${i}.ping`)}
                aria-invalid={errors?.ping ? "true" : "false"}
                error={errors?.ping?.message}
            />
        </Fragment>
    );
}

function VolumeInputs({ register, errors, i }) {
    return (
        <div className={styles.volume}>
            <TextField
                id={`methods.docker.volumes.${i}.key`}
                label="Name of the volume"
                {...register(`methods.docker.volumes.${i}.key`)}
                aria-invalid={errors?.key ? "true" : "false"}
                error={errors?.key?.message}
                placeholder="data"
                required
            />
            <TextField
                id={`methods.docker.volumes.${i}.value`}
                label="Path in the container"
                {...register(`methods.docker.volumes.${i}.value`)}
                aria-invalid={errors?.value ? "true" : "false"}
                error={errors?.value?.message}
                placeholder="/var/lib/data"
                required
            />
        </div>
    );
}

function PortInputs({ register, errors, i }) {
    return (
        <div className={styles.port}>
            <TextField
                id={`methods.docker.ports.${i}.key`}
                label="Port in the container"
                {...register(`methods.docker.ports.${i}.key`)}
                aria-invalid={errors?.key ? "true" : "false"}
                error={errors?.key?.message}
                placeholder="3000"
                required
            />
            <TextField
                id={`methods.docker.ports.${i}.value`}
                label="Port out of the container"
                {...register(`methods.docker.ports.${i}.value`)}
                aria-invalid={errors?.value ? "true" : "false"}
                error={errors?.value?.message}
                placeholder="3000"
                required
            />
        </div>
    );
}

function ContainerEnvironmentInputs({ register, errors, i }) {
    return (
        <div className={styles.containerEnvironment}>
            <TextField
                id={`methods.docker.environment.${i}.key`}
                label="Name in Docker"
                {...register(`methods.docker.environment.${i}.key`)}
                aria-invalid={errors?.key ? "true" : "false"}
                error={errors?.key?.message}
                placeholder="TOKEN"
                required
            />
            <TextField
                id={`methods.docker.environment.${i}.value`}
                label="Name in Vertex"
                {...register(`methods.docker.environment.${i}.value`)}
                aria-invalid={errors?.value ? "true" : "false"}
                error={errors?.value?.message}
                placeholder="TOKEN"
                required
            />
        </div>
    );
}

const environmentTypes = [
    { value: "string", label: "String" },
    { value: "port", label: "Port" },
    { value: "timezone", label: "Timezone" },
    { value: "url", label: "URL" },
];

const urlKinds = [
    { value: "client", label: "Client" },
    { value: "server", label: "Server" },
];

const schema = object({
    id: yup
        .string()
        .label("Service ID")
        .required()
        .matches(
            /^[a-z0-9-]+$/,
            "Can only have lowercase letters, numbers and dashes."
        ),
    name: yup.string().label("Service name").required(),
    repository: yup.string().label("Repository").url(),
    description: yup.string().label("Description").required().max(100),
    color: yup
        .string()
        .label("Color")
        .matches(/^#[0-9a-f]{6}$/i, {
            message: "Must be a valid hex color.",
            excludeEmptyString: true,
        }),

    environment: yup.array().of(
        yup.object().shape({
            type: yup
                .string()
                .label("Type")
                .oneOf(environmentTypes.map((e) => e.value))
                .required(),
            name: yup
                .string()
                .label("Name")
                .matches(
                    /^[A-Z0-9_]+$/,
                    "Can only have uppercase letters, numbers and underscores."
                )
                .required(),
            display_name: yup.string().label("Display name").required(),
            default: yup.string().label("Default value"),
            description: yup.string().label("Description").max(100),
        })
    ),

    urls: yup.array().of(
        yup.object().shape({
            name: yup.string().label("Name").required(),
            port: yup.number().label("Port").required(),
            ping: yup.string().label("Ping"),
            kind: yup
                .string()
                .label("Kind")
                .oneOf(urlKinds.map((e) => e.value))
                .required(),
        })
    ),

    methods: yup.object().shape({
        docker: yup.object().shape({
            image: yup.string().label("Docker image").required(),
            command: yup.string().label("Command"),
            volumes: yup.array().of(
                yup.object().shape({
                    key: yup.string().label("Name").required(),
                    value: yup.string().label("Path").required(),
                })
            ),
            ports: yup.array().of(
                yup.object().shape({
                    key: yup.number().label("Port input").required(),
                    value: yup.number().label("Port output").required(),
                })
            ),
            environment: yup.array().of(
                yup.object().shape({
                    key: yup.string().label("Name").required(),
                    value: yup.string().label("Name").required(),
                })
            ),
        }),
    }),
});

type FormData = yup.InferType<typeof schema>;

export default function ServiceEditor() {
    useTitle("Service Editor");

    const resolver = yupResolver(schema);
    const {
        control,
        register,
        handleSubmit,
        getValues: getService,
        formState: { errors },
    } = useForm<FormData>({ resolver });
    const onSubmit: SubmitHandler<FormData> = (data) => console.log(data);

    const {
        fields: environment,
        append: appendEnvironment,
        remove: removeEnvironment,
    } = useFieldArray({
        control,
        name: "environment",
    });

    const {
        fields: urls,
        append: appendUrl,
        remove: removeUrl,
    } = useFieldArray({
        control,
        name: "urls",
    });

    const {
        fields: volumes,
        append: appendVolume,
        remove: removeVolume,
    } = useFieldArray({
        control,
        name: "methods.docker.volumes",
    });

    const {
        fields: ports,
        append: appendPort,
        remove: removePort,
    } = useFieldArray({
        control,
        name: "methods.docker.ports",
    });

    const {
        fields: containerEnvironment,
        append: appendContainerEnvironment,
        remove: removeContainerEnvironment,
    } = useFieldArray({
        control,
        name: "methods.docker.environment",
    });

    const [yaml, setYaml] = useState();

    const download = () => {
        const service = getService();
        console.log(service);

        const s: any = produce(service, (draft: any) => {
            const volumes = draft.methods.docker.volumes;
            const newVolumes = Object.assign(
                {},
                ...volumes.map((x: any) => ({ [x.key]: x.value }))
            );

            const ports = draft.methods.docker.ports;
            const newPorts = Object.assign(
                {},
                ...ports.map((x: any) => ({ [x.key]: x.value }))
            );

            const environment = draft.methods.docker.environment;
            const newEnvironment = Object.assign(
                {},
                ...environment.map((x: any) => ({ [x.key]: x.value }))
            );

            return {
                ...draft,
                methods: {
                    docker: {
                        volumes: newVolumes,
                        ports: newPorts,
                        environment: newEnvironment,
                    },
                },
            };
        });

        console.log(s);

        api.vxServiceEditor.editor.toYaml(s).then((data) => setYaml(data));
    };

    return (
        <Content className={styles.content}>
            <Title variant="h2">Info</Title>
            <div className={styles.inputs}>
                <TextField
                    id="id"
                    label="Service ID"
                    {...register("id", { required: true })}
                    aria-invalid={errors.id ? "true" : "false"}
                    error={errors.id?.message}
                    placeholder="my-service"
                    description="Lowercase identifier for the service."
                    required
                />
                <TextField
                    id="name"
                    label="Service name"
                    {...register("name")}
                    aria-invalid={errors.name ? "true" : "false"}
                    error={errors.name?.message}
                    placeholder="My service"
                    description="Human-readable name for the service."
                    required
                />
                <TextField
                    id="repository"
                    label="Repository"
                    type="url"
                    {...register("repository")}
                    aria-invalid={errors.repository ? "true" : "false"}
                    error={errors.repository?.message}
                    placeholder="https://github.com/username/repo"
                    description="URL of the repository."
                />
                <TextField
                    id="description"
                    label="Description"
                    {...register("description")}
                    aria-invalid={errors.description ? "true" : "false"}
                    error={errors.description?.message}
                    placeholder="A simple database watcher."
                    description="Short description of the service."
                    required
                />
                <TextField
                    id="color"
                    label="Color"
                    {...register("color")}
                    aria-invalid={errors.color ? "true" : "false"}
                    error={errors.color?.message}
                    placeholder="#d73d3d"
                    description="Color of the service used in the UI."
                />
            </div>

            <Horizontal className={styles.title} alignItems="center">
                <Title variant="h2">Environment</Title>
                <Spacer />
                <Button
                    variant="colored"
                    rightIcon={<MaterialIcon icon="add" />}
                    onClick={() =>
                        appendEnvironment({
                            type: "string",
                        })
                    }
                >
                    Add variable
                </Button>
            </Horizontal>
            {environment.length === 0 && (
                <NoItems icon="list" text="No environment variables" />
            )}
            <List>
                {environment.map((field, i) => {
                    return (
                        <ListItem key={field.id}>
                            <Vertical
                                alignItems="stretch"
                                style={{ width: "100%" }}
                                gap={5}
                            >
                                <Horizontal>
                                    <Title variant="h3">
                                        Environment variable #{i + 1}
                                    </Title>
                                    <Spacer />
                                    <Button
                                        variant="danger"
                                        rightIcon={
                                            <MaterialIcon icon="delete" />
                                        }
                                        onClick={() => removeEnvironment(i)}
                                    >
                                        Remove
                                    </Button>
                                </Horizontal>
                                <div className={styles.inputs}>
                                    <EnvironmentInputs
                                        key={field.id}
                                        control={control}
                                        register={register}
                                        i={i}
                                        errors={errors.environment?.[i]}
                                    />
                                </div>
                            </Vertical>
                        </ListItem>
                    );
                })}
            </List>

            <Horizontal className={styles.title} alignItems="center">
                <Title variant="h2">URLs</Title>
                <Spacer />
                <Button
                    variant="colored"
                    rightIcon={<MaterialIcon icon="add" />}
                    onClick={() => appendUrl({ kind: "client" })}
                >
                    Add URL
                </Button>
            </Horizontal>
            {urls.length === 0 && <NoItems icon="public" text="No URLs." />}
            <List>
                {urls.map((field, i) => {
                    return (
                        <ListItem key={field.id}>
                            <Vertical
                                alignItems="stretch"
                                style={{ width: "100%" }}
                                gap={5}
                            >
                                <Horizontal>
                                    <Title variant="h3">URL #{i + 1}</Title>
                                    <Spacer />
                                    <Button
                                        variant="danger"
                                        rightIcon={
                                            <MaterialIcon icon="delete" />
                                        }
                                        onClick={() => removeUrl(i)}
                                    >
                                        Remove
                                    </Button>
                                </Horizontal>
                                <div className={styles.inputs}>
                                    <UrlInputs
                                        key={field.id}
                                        control={control}
                                        register={register}
                                        i={i}
                                        errors={errors.urls?.[i]}
                                    />
                                </div>
                            </Vertical>
                        </ListItem>
                    );
                })}
            </List>

            <Title variant="h2">Docker</Title>
            <Vertical className={classNames(styles.inputs)} gap={15}>
                <TextField
                    id="methods.docker.image"
                    label="Docker image"
                    {...register("methods.docker.image")}
                    aria-invalid={
                        errors.methods?.docker?.image ? "true" : "false"
                    }
                    error={errors.methods?.docker?.image?.message}
                    placeholder="org/repo, ghcr.io/org/repo, etc."
                    description="The image to pull from a registry."
                    required
                />
                <TextField
                    id="methods.docker.command"
                    label="Command"
                    {...register("methods.docker.command")}
                    aria-invalid={
                        errors.methods?.docker?.command ? "true" : "false"
                    }
                    error={errors.methods?.docker?.command?.message}
                    placeholder="npm start"
                    description="A command to run on startup."
                />
            </Vertical>

            <Horizontal className={styles.title} alignItems="center">
                <Title variant="h2">Docker Volumes</Title>
                <Spacer />
                <Button
                    variant="colored"
                    rightIcon={<MaterialIcon icon="add" />}
                    onClick={() => appendVolume({})}
                >
                    Add volume
                </Button>
            </Horizontal>
            {volumes.length === 0 ? (
                <NoItems icon="storage" text="No Docker volumes." />
            ) : (
                <Card>
                    {volumes.map((field, i) => {
                        return (
                            <div
                                key={field.id}
                                className={classNames({
                                    [styles.inputsRow]: true,
                                    [styles.inputsRowFirst]: i === 0,
                                })}
                            >
                                <VolumeInputs
                                    key={field.id}
                                    register={register}
                                    i={i}
                                    errors={
                                        errors?.methods?.docker?.volumes?.[i]
                                    }
                                />
                                <Button
                                    variant="danger"
                                    className={i === 0 && styles.deleteOffset}
                                    rightIcon={<MaterialIcon icon="delete" />}
                                    onClick={() => removeVolume(i)}
                                >
                                    Remove
                                </Button>
                            </div>
                        );
                    })}
                </Card>
            )}

            <Horizontal className={styles.title} alignItems="center">
                <Title variant="h2">Docker Ports</Title>
                <Spacer />
                <Button
                    variant="colored"
                    rightIcon={<MaterialIcon icon="add" />}
                    onClick={() => appendPort({})}
                >
                    Add port
                </Button>
            </Horizontal>
            {ports.length === 0 ? (
                <NoItems icon="hub" text="No Docker ports." />
            ) : (
                <Card>
                    {ports.map((field, i) => {
                        return (
                            <div
                                key={field.id}
                                className={classNames({
                                    [styles.inputsRow]: true,
                                    [styles.inputsRowFirst]: i === 0,
                                })}
                            >
                                <PortInputs
                                    key={field.id}
                                    register={register}
                                    i={i}
                                    errors={errors?.methods?.docker?.ports?.[i]}
                                />
                                <Button
                                    variant="danger"
                                    className={i === 0 && styles.deleteOffset}
                                    rightIcon={<MaterialIcon icon="delete" />}
                                    onClick={() => removePort(i)}
                                >
                                    Remove
                                </Button>
                            </div>
                        );
                    })}
                </Card>
            )}

            <Horizontal className={styles.title} alignItems="center">
                <Title variant="h2">Docker Environments</Title>
                <Spacer />
                <Button
                    variant="colored"
                    rightIcon={<MaterialIcon icon="add" />}
                    onClick={() => appendContainerEnvironment({})}
                >
                    Add port
                </Button>
            </Horizontal>
            {containerEnvironment.length === 0 ? (
                <NoItems icon="list" text="No Docker environment variables." />
            ) : (
                <Card>
                    {containerEnvironment.map((field, i) => {
                        return (
                            <div
                                key={field.id}
                                className={classNames({
                                    [styles.inputsRow]: true,
                                    [styles.inputsRowFirst]: i === 0,
                                })}
                            >
                                <ContainerEnvironmentInputs
                                    key={field.id}
                                    register={register}
                                    i={i}
                                    errors={
                                        errors?.methods?.docker?.environment?.[
                                            i
                                        ]
                                    }
                                />
                                <Button
                                    variant="danger"
                                    className={i === 0 && styles.deleteOffset}
                                    rightIcon={<MaterialIcon icon="delete" />}
                                    onClick={() =>
                                        removeContainerEnvironment(i)
                                    }
                                >
                                    Remove
                                </Button>
                            </div>
                        );
                    })}
                </Card>
            )}

            <Title variant="h2">Service.yml</Title>
            <Code language="yaml">{yaml}</Code>

            <Horizontal gap={10}>
                <Spacer />
                <Button
                    rightIcon={<MaterialIcon icon="check" />}
                    onClick={handleSubmit(onSubmit)}
                >
                    Validate
                </Button>
                <Button
                    variant="colored"
                    rightIcon={<MaterialIcon icon="download" />}
                    onClick={download}
                >
                    service.yml
                </Button>
            </Horizontal>
        </Content>
    );
}
