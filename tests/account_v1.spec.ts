import { test, expect } from "@playwright/test";
import { loginAccount, resetCache, resetDB, seedAccounts } from "./factory";
import accounts from "./accounts.json";

test.beforeAll(async ({ request }) => {
  await resetDB(request);
  await resetCache(request);
  await seedAccounts(request);
});

test.describe("Create Account", () => {
  test("creating an admin account with insufficient permissions fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.ziemowit);

    const resp = await request.post("/v1/account/admin", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Jack",
          username: "jackie420",
          password: "topsecret",
        },
      }),
    });

    expect(resp.status()).toBe(401);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("permission-denied");
  });

  test("creating a secritary account with insufficient permissions fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.ziemowit);

    const resp = await request.post("/v1/account/secritary", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Jack",
          username: "jackie420",
          password: "topsecret",
        },
      }),
    });

    expect(resp.status()).toBe(401);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("permission-denied");
  });

  test("creating an admin account with empty username fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/account/admin", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Jack",
          username: "",
          password: "topsecret",
        },
      }),
    });

    expect(resp.status()).toBe(400);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("invalid-account-username");
  });

  test("creating a secritary account with empty username fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/account/secritary", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Jack",
          username: "",
          password: "topsecret",
        },
      }),
    });

    expect(resp.status()).toBe(400);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("invalid-account-username");
  });

  test("creating an admin account with empty password fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/account/admin", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Jack",
          username: "jackie420",
          password: "",
        },
      }),
    });

    expect(resp.status()).toBe(400);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("invalid-account-password");
  });

  test("creating a secritary account with empty password fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/account/secritary", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Jack",
          username: "jackie420",
          password: "",
        },
      }),
    });

    expect(resp.status()).toBe(400);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("invalid-account-password");
  });

  test("creating an admin account with empty display name fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/account/admin", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "",
          username: "jackie420",
          password: "topsecret",
        },
      }),
    });

    expect(resp.status()).toBe(400);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("invalid-account-display-name");
  });

  test("creating a secritary account with empty display name fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/account/secritary", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "",
          username: "jackie420",
          password: "topsecret",
        },
      }),
    });

    expect(resp.status()).toBe(400);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("invalid-account-display-name");
  });

  test("creating an admin account with an existing username fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/account/admin", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: accounts.harvey,
      }),
    });

    expect(resp.status()).toBe(409);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("account-exists");
  });

  test("creating a secritary account with an existing username fails", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/account/secritary", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: accounts.cody,
      }),
    });

    expect(resp.status()).toBe(409);

    const respBody = await resp.json();
    expect(respBody.error_id).toBe("account-exists");
  });

  test("creating an admin account with valid parameters and permissions works", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/account/admin", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Jack",
          username: "jackie420",
          password: "topsecret",
        },
      }),
    });

    expect(resp).toBeOK();
  });

  test("creating a secritary account with valid parameters and permissions works", async ({
    request,
  }) => {
    const token = await loginAccount(request, accounts.rex);

    const resp = await request.post("/v1/account/secritary", {
      headers: {
        Authorization: token,
      },
      data: JSON.stringify({
        new_account: {
          display_name: "Dank",
          username: "dankie444",
          password: "topsecret",
        },
      }),
    });

    expect(resp).toBeOK();
  });
});
