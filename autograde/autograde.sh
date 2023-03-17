#!/bin/bash -p

if [ "$#" -ne 2 ]; then
	echo "Usage: $0 <assignment> <student-username>" >&2
	exit 1
fi

ASSIGNMENT="$1"
STUDENT_USERNAME="$2"

umask 000

SUPPORT_CODE=/course/cs666/autograde/stencils/$ASSIGNMENT

# Remove any stencil files the student submitted
# (we'll replace with our own versions from the support code)
for f in $SUPPORT_CODE/*; do
    echo "Removing stencil file $(basename $f)"
    rm -f ./$(basename $f)
done

# figure out what .go files there are before
# copying in stencil code so that we don't
# blocklist imports in our own code
GO_FILES=$(echo *.go)

if [ "$GO_FILES" == "*.go" ]; then
	echo "No .go files given" >&2
	exit 1
fi

# copy framework code into student submission dir
cp $SUPPORT_CODE/* .

# Remove potentially-dangerous imports
/course/cs666/tabin/blocklist_imports $GO_FILES || exit 1

# recalculate to include our .go files this time
GO_FILES=$(echo *.go)

# First attempt the build to see if it fails
TEMPLATE=/tmp/cs666_bin_tmp.XXXXXXXXXX
TMP=$(mktemp $TEMPLATE) || exit 1
go build -o $TMP $GO_FILES || exit 1

# It succeeded; move it to a permanent location
PERM=/tmp/cs666_bin
cat $TMP > $PERM
chmod u=rwx,g=rx,o= $PERM
rm $TMP

function ingest {
	echo "Executing:"
	while read line; do
		# lines are of the form:
		# problem:<name> points:<points> comment:<comment>
		PROB=$(echo $line | cut -d ' ' -f 1 | cut -d : -f 2)
		POINTS=$(echo $line | cut -d ' ' -f 2 | cut -d : -f 2)
		COMMENT=$(echo $line | cut -d ' ' -f 3- | cut -d : -f 2-)
		echo "modifydb --command grade -p \"$ASSIGNMENT:$PROB\" --points \"$POINTS\" "\
			"--comment \"$COMMENT\" -s \"$STUDENT_USERNAME\" --files $GO_FILES"
		/course/cs666/tabin/modifydb --command grade -p "$ASSIGNMENT:$PROB" --points "$POINTS" \
			--comment "$COMMENT" -s "$STUDENT_USERNAME" --files $GO_FILES
	done
}

$PERM | ingest

rm $PERM

# Now check for any solution files
if ! [ -d /course/cs666/secret/$STUDENT_USERNAME/$ASSIGNMENT ]; then
	# no solution files to check; we're done
	exit 0
fi

# loop over the files one at a time, and make sure
# that the handed in file matches the reference
# file exactly
ls /course/cs666/secret/$STUDENT_USERNAME/$ASSIGNMENT | while read file; do
	if ! [ -f $file ]; then
		echo "Handin does not include required file $file" >&2
		# ingest uses the assignment $ASSIGNMENT:$PROB, so this will
		# become $ASSIGNMENT:$file, which works
		echo "prob:$file points:0 comment:file not present" | ingest
	else
		SAME="false"
		diff /course/cs666/secret/$STUDENT_USERNAME/$ASSIGNMENT/$file $file >/dev/null 2>/dev/null
		if [ $? -eq 0 ]; then
			SAME="true"
		fi
		if [ "$SAME" == "true" ]; then
			# ingest uses the assignment $ASSIGNMENT:$PROB, so this will
			# become $ASSIGNMENT:$file, which works
			echo "prob:$file points:100 comment:files match" | ingest
		else
			# ingest uses the assignment $ASSIGNMENT:$PROB, so this will
			# become $ASSIGNMENT:$file, which works
			echo "prob:$file points:0 comment:files do not match" | ingest
		fi
	fi
done
