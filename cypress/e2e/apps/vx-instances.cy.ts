const createContainerAndNavigate = (to: string) => {
    cy.visit("/app/containers");
    cy.request(
        "POST",
        "http://localhost:6130/api/service/postgres/install"
    ).then((res: any) => {
        cy.visit(`/app/containers/${res.body.uuid}${to}`);
    });
};

describe("The Vertex Containers app", () => {
    it("loads", () => {
        cy.visit("/app/containers");
    });

    it("can create a new container", () => {
        // Navigate to the create container page
        cy.visit("/app/containers");
        cy.contains("Create container").click();
        cy.url().should("include", "/app/containers/add");

        // Create an container
        cy.contains("Postgres").click();
        cy.get("button").contains("Create container").click();

        // Go back
        cy.visit("/app/containers");
        cy.contains("Postgres");
    });

    it("can navigate to an container", () => {
        cy.request(
            "POST",
            "http://localhost:6130/api/app/containers/service/postgres/install"
        );

        // Navigate to
        cy.visit("/app/containers");
        cy.contains("Postgres").click();
    });

    it("can navigate to the container home page", () => {
        createContainerAndNavigate("/home");
        cy.get("h2").contains("URLs");
    });

    it("can navigate to the container logs page", () => {
        createContainerAndNavigate("/logs");
        cy.get("h2").contains("Logs");
    });

    it("can navigate to the container Docker page", () => {
        createContainerAndNavigate("/docker");
        cy.get("h2").contains("Container");
        cy.get("h2").contains("Image");
    });

    it("can navigate to the container environment page", () => {
        createContainerAndNavigate("/environment");
        cy.get("h2").contains("Environment");
    });

    it("can navigate to the container updates page", () => {
        createContainerAndNavigate("/update");
        cy.get("h2").contains("Update");
    });

    it("can navigate to the container settings page", () => {
        createContainerAndNavigate("/settings");
        cy.get("h2").contains("Settings");
    });

    it("can delete an container", () => {
        // Go to the containers page
        cy.visit("/app/containers");

        // Create an container
        createContainerAndNavigate("/");

        // Delete the container
        cy.contains("Delete").click();
        cy.get("button").contains("Confirm").click();
        cy.get("button").contains("Delete").should("not.exist");

        // The url should not end with the container uuid
        cy.url().should("not.match", /\/[a-z0-9-]{36}$/);
    });
});
