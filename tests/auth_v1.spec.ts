import { test, expect } from "@playwright/test";
import { createAccountAsSuperAdmin, resetCache, resetDB } from "./factory";

test.beforeAll(async ({ request }) => {
  await resetDB(request);
  await resetCache(request);
});

test.describe("Login", () => {
  test("when user provides incorrect login credentials login fails", async ({
    request,
  }) => {
    const resp = await request.post("/v1/login/username", {
      data: JSON.stringify({
        username: "foo",
        password: "bar",
      }),
    });

    expect(resp.status()).toBe(401);
  });

  test("when user provides correct login credentials login fails", async ({
    request,
  }) => {
    await createAccountAsSuperAdmin(request, "ziemowit");

    const resp = await request.post("/v1/login/username", {
      data: JSON.stringify({
        username: "ziemowit",
        password: "1234",
      }),
    });

    expect(resp).toBeOK();

    const respBody = await resp.json();
    expect(respBody.session_token).not.toBeUndefined();
  });
});
