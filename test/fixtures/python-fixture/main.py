# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import asyncio
from typing import Optional
from ztvs_sdk import Metadata, run, Check, Finding # type: ignore

class PythonCheck(Check):
    def id(self) -> str:
        return "python_test_check"

    def name(self) -> str:
        return "Python Fixture Check"

    async def run(self) -> Optional[Finding]:
        return Finding(
            id="F-PY-001",
            severity="info",
            title="Python Fixture Executed",
            description="This is a successful result from the Python polyglot fixture.",
            evidence={"runtime": "python3"},
            remediation="No action required."
        )

async def main():
    meta = Metadata(
        name="python-fixture",
        version="1.0.0",
        api_version=1
    )
    await run(meta, [PythonCheck()])

if __name__ == "__main__":
    asyncio.run(main())
