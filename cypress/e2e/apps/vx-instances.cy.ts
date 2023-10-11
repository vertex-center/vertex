const createInstanceAndNavigate = (to: string) => {
    cy.visit("/app/vx-instances");
    cy.request(
        "POST",
        "http://localhost:6130/api/app/vx-instances/service/postgres/install"
    ).then((res: any) => {
        cy.visit(`/app/vx-instances/${res.body.uuid}${to}`);
    });
};

describe("The Vertex Instances app", () => {
    it("loads", () => {
        cy.visit("/app/vx-instances");
    });

    it("can create a new instance", () => {
        // Navigate to the create instance page
        cy.visit("/app/vx-instances");
        cy.contains("Create instance").click();
        cy.url().should("include", "/app/vx-instances/add");

        // Create an instance
        cy.contains("Postgres").click();
        cy.get("button").contains("Create instance").click();

        // Go back
        cy.visit("/app/vx-instances");
        cy.contains("Postgres");
    });

    it("can navigate to an instance", () => {
        cy.request(
            "POST",
            "http://localhost:6130/api/app/vx-instances/service/postgres/install"
        );

        // Navigate to
        cy.visit("/app/vx-instances");
        cy.contains("Postgres").click();
    });

    it("can navigate to the instance home page", () => {
        createInstanceAndNavigate("/home");
        cy.get("h2").contains("URLs");
    });

    it("can navigate to the instance logs page", () => {
        createInstanceAndNavigate("/logs");
        cy.get("h2").contains("Logs");
    });

    it("can navigate to the instance Docker page", () => {
        createInstanceAndNavigate("/docker");
        cy.get("h2").contains("Container");
        cy.get("h2").contains("Image");
    });

    it("can navigate to the instance environment page", () => {
        createInstanceAndNavigate("/environment");
        cy.get("h2").contains("Environment");
    });

    it("can navigate to the instance updates page", () => {
        createInstanceAndNavigate("/update");
        cy.get("h2").contains("Update");
    });

    it("can navigate to the instance settings page", () => {
        createInstanceAndNavigate("/settings");
        cy.get("h2").contains("Settings");
    });

    it("can delete an instance", () => {
        // Go to the instances page
        cy.visit("/app/vx-instances");

        // Create an instance
        createInstanceAndNavigate("/");

        // Delete the instance
        cy.contains("Delete").click();
        cy.get("button").contains("Confirm").click();
        cy.get("button").contains("Delete").should("not.exist");

        // The url should not end with the instance uuid
        cy.url().should("not.match", /\/[a-z0-9-]{36}$/);
    });
});
