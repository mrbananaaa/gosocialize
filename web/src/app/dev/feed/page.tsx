import Sidebar from "@/components/ui/sidebar";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Gosocialize",
  description: "I dunno what tf was ts.",
};

export default function FeedPage() {
  return (
    <div className="relative w-full min-h-screen bg-[#070708] overflow-hidden font-sans antialiased">
      <div className="absolute inset-0 z-0 pointer-events-none overflow-hidden">
        <div className="absolute bottom-[0%] right-[0%] w-[70vw] h-[70vw] max-w-[750px] max-h-[750px] rounded-full bg-[#33312E] opacity-25 blur-[140px]" />
        <div className="absolute bottom-[4%] right-[4%] w-[40vw] h-[40vw] max-w-[480px] max-h-[480px] rounded-full bg-[#F5F2EB] opacity-18 blur-[100px]" />
      </div>

      <div className="relative min-h-screen flex justify-center w-full">
        <div className="flex w-full max-w-7xl">
          <aside className="hidden sm:flex w-20 xl:w-64 h-screen sticky top-0 p-3 border-r border-gray-800">
            <Sidebar />
          </aside>

          <main className="min-h-screen w-full p-3">main page</main>
        </div>
      </div>
    </div>
  );
}
