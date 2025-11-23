import MintProposals from "./mint-proposals";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import AddOrUpdateAgentImageProposals from "./add-or-update-agent-image-proposals";

export default function DashboardModeration() {
  return (
    <div className="flex-1 space-y-4 p-8 pt-6">
      <div className="flex items-center justify-between space-y-2">
        <h2 className="text-3xl font-bold tracking-tight">Moderation</h2>
      </div>
      <div>
        <Tabs defaultValue="mint-proposals" className="gap-6">
          <TabsList className="gap-x-2">
            <TabsTrigger className="" value="mint-proposals">
              Mint proposals
            </TabsTrigger>
            <TabsTrigger className="" value="agent-image-proposals">
              Agent images proposals
            </TabsTrigger>
          </TabsList>
          <TabsContent value="mint-proposals">
            <MintProposals />
          </TabsContent>
          <TabsContent value="agent-image-proposals">
            <AddOrUpdateAgentImageProposals />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
