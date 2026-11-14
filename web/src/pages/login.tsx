import { LoginForm } from "@/components/login/login-form"

export default function LoginPage() {
  return (
    <div className="grid min-h-svh lg:grid-cols-2">
      <div className="relative hidden bg-zinc-900 text-white lg:block">
        <div className="flex h-full flex-col justify-between p-6 md:p-10">
          <div className="absolute inset-0 h-full w-full p-30">
            <img
              src="/logo-1024x1024.png"
              alt="Logo"
              className="h-full w-full object-contain"
            />
          </div>
          <div className="flex text-lg font-normal">Agentic AI Platform</div>
          <div className="flex text-lg font-extralight">
            Decentralized Intelligence at Your Fingertips
          </div>
        </div>
      </div>
      <div className="flex flex-col gap-4 p-6 md:p-10">
        <div className="flex justify-center gap-2 md:justify-end">
          <div className="flex items-center gap-2 font-medium">
            Swiss 6022.
            <div className="flex h-6 w-6 items-center justify-center rounded-md bg-primary text-primary-foreground">
              <img src="/logo-192x192.png" alt="Logo" className="h-4 w-4" />
            </div>
          </div>
        </div>
        <div className="flex flex-1 items-center justify-center">
          <div className="w-full max-w-xs">
            <p className="mb-6 text-center text-lg font-medium lg:hidden">
              Agentic AI Platform
            </p>
            <LoginForm />
          </div>
        </div>
      </div>
    </div>
  )
}
