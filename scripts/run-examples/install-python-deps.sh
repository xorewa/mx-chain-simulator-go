#!/bin/bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repository_root="$(cd "${script_dir}/../.." && pwd)"
venv_dir="${PYTHON_VENV_DIR:-${script_dir}/.venv}"

# Keep example dependencies isolated from the host Python installation. This
# avoids requiring sudo and works on distributions that protect the system
# interpreter (PEP 668).
python3 -m venv "${venv_dir}"
"${venv_dir}/bin/python" -m pip install --upgrade pip
"${venv_dir}/bin/python" -m pip install -r "${repository_root}/examples/requirements.txt" pytest
