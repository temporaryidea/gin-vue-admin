import json


def validate_message():
    # Read test_message_v3.py
    with open('test_message_v3.py', 'r') as f:
        content = f.read()

    # Extract message content
    message_start = content.find('message = json.dumps({')
    message_end = content.find('    })', message_start) + 6
    message_str = content[message_start:message_end]

    # Execute the message assignment
    local_dict = {}
    exec(message_str, {'json': json}, local_dict)
    message = local_dict['message']

    # Perform validations
    message_size = len(message)
    print(f'Message size (bytes): {message_size}')
    size_status = "OK" if message_size < 150 * 1024 else "Warning - Too large"
    print(f'Size limit check: {size_status}')

    # Check for double encoding
    try:
        decoded = json.loads(message)
        try:
            json.loads(decoded)
            print('Double encoding check: Warning - Double encoded')
        except json.JSONDecodeError:
            print('Double encoding check: OK - Single encoded')
    except json.JSONDecodeError:
        print('JSON validation: Invalid JSON')
        return

    print('JSON validation: Valid')

    # Check msg_type
    if 'msg_type="interactive"' in content:
        print('msg_type check: OK - interactive')
    else:
        print('msg_type check: Warning - msg_type missing or incorrect')


if __name__ == '__main__':
    validate_message()
