import Layout from "@/components/Layout"
import { NextPageWithLayout } from "./_app"
import { ReactElement, useEffect } from "react"
import { AreaChart, ResponsiveContainer, XAxis, YAxis, Area, Legend, Tooltip } from "recharts"

const testData = [
    {
        date: "2021-01-01",
        visitors: 100
    },
    {
        date: "2021-01-02",
        visitors: 200
    },
    {
        date: "2021-01-03",
        visitors: 300
    },
    {
        date: "2021-01-04",
        visitors: 400
    },
    {
        date: "2021-01-05",
        visitors: 200
    },
    {
        date: "2021-01-06",
        visitors: 250
    },
    {
        date: "2021-01-07",
        visitors: 300
    },
]

const Page: NextPageWithLayout = () => {
    useEffect(() => {
        console.log(window.sessionStorage)
    }, [])
    return (
      <main className="p-12">
        <div className="container mx-auto h-12">
            <div className="columns-1 p-3">
                <div className="container p-4 bg-gray-100 mx-auto rounded">
                    <div className="bg-white rounded shadow p-6 mb-0 pb-0">
                        <ResponsiveContainer width="100%" height={400}>
                            <AreaChart
                                data={testData}
                            >
                                <XAxis dataKey="date"/>
                                <YAxis/>
                                {/* <Legend/> */}
                                <Tooltip/>
                                <Area type="monotone" dataKey="visitors" stroke="rgb(147 197 253)" fill="rgb(147 197 253)"/>
                            </AreaChart>
                        </ResponsiveContainer>
                    </div>
                </div>
            </div>
        </div>
      </main>
    )
}


Page.getLayout = function getLayout(page: ReactElement) {
    return (
        <Layout page="metrics">
            {page}
        </Layout>
    )
}

export default Page