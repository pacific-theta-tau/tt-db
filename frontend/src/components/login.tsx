import React from 'react';

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"


const Login: React.FC = () => {
  return (
    <Card className="w-[400px]">
      <CardHeader>
        <CardTitle>Login</CardTitle>
        <CardDescription>Enter your username and password to access the database</CardDescription>
      </CardHeader>
      <CardContent>
        <form>
            <div className="flex flex-col space-y-1.5">
              <Label htmlFor="name">Name</Label>
              <Input id="name" placeholder="Name of your project" />
            </div>
            <div className="flex flex-col space-y-1.5">
                <Label htmlFor="password">Password</Label>
                <Input id="password" />
            </div>
        </form>
      </CardContent>
      <CardFooter className="content-center">
        <Button>Deploy</Button>
      </CardFooter>
    </Card>
  )
}

export default Login
