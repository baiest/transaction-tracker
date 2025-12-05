// middleware.ts
import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export function middleware(req: NextRequest) {
  const token = req.cookies.get("token")?.value;

  if (!token)
    return NextResponse.redirect(
      process.env.LOGIN_URL! || "https://google.com"
    );

  // const validate = await fetch(`${process.env.API_URL}/auth/validate`, {
  //   headers: {
  //     Authorization: `Bearer ${token}`
  //   }
  // });

  // if (validate.status === 401) {
  //   const response = NextResponse.redirect(new URL("/login", req.url));

  //   response.cookies.set({
  //     name: "token",
  //     value: "",
  //     maxAge: 0
  //   });

  //   return response;
  // }

  return NextResponse.next();
}

export const config = {
  matcher: ["/", "/movements"]
};
